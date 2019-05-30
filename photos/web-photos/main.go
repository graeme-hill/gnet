package main

func main() {
	ctx := context.Background()
	errChan := server.Run(ctx, server.Options{
		Addr: ":8000",
		EventStoreConnStr: ":memory:",
		FileStoreConnStr: ":memory:",
	})

	select {
	case err := <-errChan:
		if err != nil {
			log.Printf("Server stopped with error: %v", err)
		} else {
			log.Printf("Server stopped cleanly")
		}
	}
}