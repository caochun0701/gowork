package memqueue

import (
	"libbeat/beat"
	"libbeat/common/atomic"
	"libbeat/publisher"
	"libbeat/publisher/queue"
)

type forgetfullProducer struct {
	broker    *Broker
	openState openState
}

type ackProducer struct {
	broker    *Broker
	cancel    bool
	seq       uint32
	state     produceState
	openState openState
}

type openState struct {
	log    logger
	isOpen atomic.Bool
	done   chan struct{}
	events chan pushRequest
}

type produceState struct {
	cb        ackHandler
	dropCB    func(beat.Event)
	cancelled bool
	lastACK   uint32
}

type ackHandler func(count int)

func newProducer(b *Broker, cb ackHandler, dropCB func(beat.Event), dropOnCancel bool) queue.Producer {
	openState := openState{
		log:    b.logger,
		isOpen: atomic.MakeBool(true),
		done:   make(chan struct{}),
		events: b.events,
	}

	if cb != nil {
		p := &ackProducer{broker: b, seq: 1, cancel: dropOnCancel, openState: openState}
		p.state.cb = cb
		p.state.dropCB = dropCB
		return p
	}
	return &forgetfullProducer{broker: b, openState: openState}
}

func (p *forgetfullProducer) Publish(event publisher.Event) bool {
	return p.openState.publish(p.makeRequest(event))
}

func (p *forgetfullProducer) TryPublish(event publisher.Event) bool {
	return p.openState.tryPublish(p.makeRequest(event))
}

func (p *forgetfullProducer) makeRequest(event publisher.Event) pushRequest {
	return pushRequest{event: event}
}

func (p *forgetfullProducer) Cancel() int {
	p.openState.Close()
	return 0
}

func (p *ackProducer) Publish(event publisher.Event) bool {
	return p.updSeq(p.openState.publish(p.makeRequest(event)))
}

func (p *ackProducer) TryPublish(event publisher.Event) bool {
	return p.updSeq(p.openState.tryPublish(p.makeRequest(event)))
}

func (p *ackProducer) updSeq(ok bool) bool {
	if ok {
		p.seq++
	}
	return ok
}

func (p *ackProducer) makeRequest(event publisher.Event) pushRequest {
	req := pushRequest{
		event: event,
		seq:   p.seq,
		state: &p.state,
	}
	return req
}

func (p *ackProducer) Cancel() int {
	p.openState.Close()

	if p.cancel {
		ch := make(chan producerCancelResponse)
		p.broker.pubCancel <- producerCancelRequest{
			state: &p.state,
			resp:  ch,
		}

		// wait for cancel to being processed
		resp := <-ch
		return resp.removed
	}
	return 0
}

func (st *openState) Close() {
	st.isOpen.Store(false)
	close(st.done)
}

func (st *openState) publish(req pushRequest) bool {
	select {
	case st.events <- req:
		return true
	case <-st.done:
		st.events = nil
		return false
	}
}

func (st *openState) tryPublish(req pushRequest) bool {
	select {
	case st.events <- req:
		return true
	case <-st.done:
		st.events = nil
		return false
	default:
		st.log.Debugf("Dropping event, queue is blocked (seq=%v) ", req.seq)
		return false
	}
}