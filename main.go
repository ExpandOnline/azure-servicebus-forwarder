package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	servicebus "github.com/Azure/azure-service-bus-go"
)

type MessageHandler struct {
	Publisher *Publisher
}

func (mp MessageHandler) Handle(ctx context.Context, msg *servicebus.Message) error {

	err := mp.Publisher.Publish(msg)

	fmt.Println("Got message")

	if err != nil {
		return msg.Abandon(ctx)
	}

	return msg.Complete(ctx)
}

type Output struct {
	Info  *log.Logger
	Error *log.Logger
}

func main() {

	sbTopicName := os.Getenv("SB_TOPIC")               //"serviceproject-created-topic"
	sbSubscriptionName := os.Getenv("SB_SUBSCRIPTION") //"PMPUpdater"
	sbConnStr := os.Getenv("SB_CONN")                  // Endpoint=...

	destConnStr := os.Getenv("DEST_CONN")      // Destination connection string
	destExchange := os.Getenv("DEST_EXCHANGE") // Destination exchange

	output := Output{
		Info:  log.New(os.Stdout, "["+sbTopicName+"/"+sbSubscriptionName+"] ", log.Ldate|log.Ltime),
		Error: log.New(os.Stderr, "["+sbTopicName+"/"+sbSubscriptionName+"] ", log.Ldate|log.Ltime|log.Lshortfile),
	}

	output.Info.Println("Initializing connection")

	if sbConnStr == "" || sbTopicName == "" || sbSubscriptionName == "" {
		output.Error.Fatalln("Invalid source: Missing connection, topic or subscription")
	}

	if destConnStr == "" || destExchange == "" {
		output.Error.Fatalln("Invalid destination: Missing connection or exchange")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	output.Info.Println("Creating namespace client...")
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(sbConnStr))
	if err != nil {
		output.Error.Fatalln(err)
	}

	startListener(ctx, ns, sbTopicName, sbSubscriptionName, &output, destConnStr, destExchange)
}

func startListener(ctx context.Context, ns *servicebus.Namespace, topic string, sub string, output *Output, destConnStr string, destExchange string) {
	listener := Listener{}

	listener.SetOutput(output)
	listener.SetTopicName(topic)
	listener.SetSubscriptionName(sub)

	pub := Publisher{}
	pub.SetOutput(output)
	pub.SetURL(destConnStr)
	pub.SetExchange(destExchange)
	pub.Init()

	defer pub.Close()
	setupCloseHandler(&pub)

	handler := MessageHandler{
		Publisher: &pub,
	}

	listener.Listen(ctx, ns, handler)
}

func setupCloseHandler(pub *Publisher) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		pub.Close()
		os.Exit(0)
	}()
}
