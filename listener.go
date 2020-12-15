package main

import (
	"context"

	servicebus "github.com/Azure/azure-service-bus-go"
)

// Listener ...
type Listener struct {
	TopicName        string
	SubscriptionName string
	IsDLQ            bool
	output           *Output
}

// SetOutput ...
func (l *Listener) SetOutput(output *Output) {
	l.output = output
}

// SetTopicName ...
func (l *Listener) SetTopicName(topic string) {
	l.TopicName = topic
}

// SetSubscriptionName ...
func (l *Listener) SetSubscriptionName(sub string) {
	l.SubscriptionName = sub
}

// SetTopicName ...
func (l *Listener) setDLQ(isDLQ bool) {
	l.IsDLQ = isDLQ
}

// Listen ...
func (l Listener) Listen(ctx context.Context, ns *servicebus.Namespace, handler MessageHandler) {

	l.output.Info.Println("Creating topic client...")
	// tm := ns.NewTopicManager()

	// topicEntity, err := l.ensureTopic(ctx, tm)

	l.output.Info.Println("Creating subscription client...")
	// sm, err := ns.NewSubscriptionManager(topicEntity.Name)
	// if err != nil {
	// 	l.output.Info.Fatalln(err)
	// }

	// subEntity, err := l.ensureSubscription(ctx, sm)

	topic, err := ns.NewTopic(l.TopicName)
	if err != nil {
		l.output.Info.Fatalln(err)
	}
	defer func() {
		_ = topic.Close(ctx)
	}()

	l.output.Info.Println("Creating subscription...")
	sub, err := topic.NewSubscription(l.SubscriptionName)
	if err != nil {
		l.output.Info.Fatalln(err)
	}

	if l.IsDLQ {
		l.output.Info.Println("Started receiving messages from DLQ..")
		for {
			l.receiveDLQ(ctx, sub, handler)
		}
	} else {
		l.output.Info.Println("Started receiving messages..")
		if err := sub.Receive(ctx, handler); err != nil {
			l.output.Info.Fatalln(err)
		}
	}

	err = sub.Close(ctx)
	if err != nil {
		l.output.Info.Fatalln(err)
	}

}

func (l Listener) receiveDLQ(ctx context.Context, sub *servicebus.Subscription, handler MessageHandler) {

	rec, err := sub.NewDeadLetterReceiver(ctx)
	if err != nil {
		l.output.Info.Fatalln(err)
	}

	if err = rec.ReceiveOne(ctx, handler); err != nil {
		l.output.Info.Fatalln(err)
	}

	rec.Close(ctx)
}

func (l Listener) ensureTopic(ctx context.Context, tm *servicebus.TopicManager, opts ...servicebus.TopicManagementOption) (*servicebus.TopicEntity, error) {

	l.output.Info.Println("Making sure topic exists...")

	name := l.TopicName
	te, err := tm.Get(ctx, name)
	if err == nil {
		l.output.Info.Println("Topic exists!")
		return te, nil
	}

	panic(err)

	// l.output.Info.Println("Topic does not exist, creating")

	// panic("Topic does not exist")

	// te, err = tm.Put(ctx, name, opts...)
	// if err != nil {
	// 	l.output.Info.Fatalln(err)
	// 	return nil, err
	// }

	// return te, nil
}

func (l Listener) ensureSubscription(ctx context.Context, sm *servicebus.SubscriptionManager, opts ...servicebus.SubscriptionManagementOption) (*servicebus.SubscriptionEntity, error) {

	l.output.Info.Println("Making sure subscription exists...")
	name := l.SubscriptionName
	subEntity, err := sm.Get(ctx, name)
	if err == nil {
		l.output.Info.Println("Subscription exists!")

		return subEntity, nil
	}

	panic(err)
	// panic("Subscription does not exist, ask Tuba!")

	// l.output.Info.Println("Subscription does not exist, creating")
	// subEntity, err = sm.Put(ctx, name, opts...)
	// if err != nil {
	// 	l.output.Info.Fatalln(err)
	// 	return nil, err
	// }

	// return subEntity, nil
}
