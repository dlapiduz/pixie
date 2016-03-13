package main

func main() {

  for {
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Enter text: ")
    text, _ := reader.ReadString('\n')
    text = strings.TrimSpace(text)

    if action := RunFilter(db, text); action.ID > 0 {
      logger.Printf("Running")
      go func() {
        out, err := RunContainer(client, action.Image)
        if err != nil {
          panic(err)
        }

        fmt.Println(out)
      }()

    }
}

