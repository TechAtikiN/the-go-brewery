# Concurrency in Golang 🔀

**Concurrency** is about **"dealing"** with many tasks at once. Go offers **built-in support** for concurrency through **goroutines** and **channels**. In this blog, let's explore how concurrency works in Go along with some common **concurrency patterns**.

## World without concurrency ☹️

Without concurrency, lines of code are executed one after another i.e **blocking** the subsequent lines until the current line is finished executing. This **synchronous execution** can be suitable for simple apps or where performance is not a concern.

However, in real-world applications, this can lead to **slow** and **unpleasant user experiences**.

Imagine Youtube app without concurrency 🚩

1. You click on a video to play, and the video starts buffering...
2. While the video is buffering, the entire UI freezes
3. You cannot scroll, pause, see comments, or interact in any way
4. Only after the video finishes loading does the app let you interact again (Boy, that would be frustrating!)

So if we are on the same boat, it's a good time to learn about **concurrency** and how it can help us build **better**, **responsive** and **fast** applications.

## Concurrency in Go 🚀

Concurrency is about "**dealing**" with many tasks at once. Strong focus on "**dealing**" with many tasks at once and not "**doing**" many tasks at once.

When tasks actually run simultaneously on different CPU cores or processors, that is called **Parallelism** (and that's a topic for another day).

**Concurrency** in Go consists of two building blocks:

- **Goroutines**
- **Channels**

### Goroutines 💪

- **Lightweight functions** that run independently.
- Created using the **`go`** keyword
- Learn more about goroutines in the previous [blog](../1-exploring-goroutines/README.md)

> #### 💡 About "go" keyword:
>
> - When you prefix a function call with the `go` keyword, it tells the Go runtime to execute that function as a goroutine. The **Go runtime** manages the **scheduling** and **execution** of goroutines on top of OS threads, due to which goroutines are **lightweight** and **efficient**.
> - Main function is also a goroutine. When the main goroutine exits, all other goroutines are terminated as well, irrespective of whether they have completed their execution or not.
> - For this to not happen, we can halt the main goroutine for some time (e.g., using time.Sleep) or use **synchronization mechanisms** like WaitGroups.
> - time.Sleep is a bad way to think about concurrency, because we don't know how long the other goroutines will take to finish, and you should almost never use it in production code for synchronization.
> - **WaitGroup** is a better way to wait for goroutines to finish.

### Channels ↔️

- Way for **goroutines** to communicate with each other
- The famous quote _**"Do not communicate by sharing memory, instead share memory by communicating"**_ holds true in Go.
- Goroutines can send and receive values through channels and they let the goroutines be aware of each other's state.
- Learn more about channels in the previous [blog](../2-exploring-channels/README.md)

## Concurrency patterns in Go 📐

As the complexity and concurrency requirements of the application increases, we need to follow certain patterns to manage the concurrency effectively and identify potential bugs.

For this, it's recommended to follow some well-known concurrency patterns.
Few of the recognized concurrency patterns in Go are~

### 1. Worker Pool Pattern

A **fixed number** of "worker" goroutines process jobs from a **shared queue**.

This pattern is essential when you need to:

- Limit concurrent operations (e.g. Database connections)
- Process a large number of tasks efficiently without overwhelming system resources

### 2. Fan-Out / Fan-In Pattern

- **Fan-Out**: Distribute a large task across multiple goroutines to perform subtasks concurrently
- **Fan-In**: Collect results from multiple goroutines into a single channel

### 3. Pipeline Pattern

In this pattern, data flows through a series of stages, each represented by a goroutine connected by channels. Each stage:

- Receives data from an input channel
- Performs its specific transformation
- Sends results to the next stage via an output channel

### 4.Generator Pattern

Used for creating streams of data that can be processed to produce some output.

- It comprises of a generator function that produces values and sends them to a channel.
- Other goroutines can receive data from this channel and process it as needed.

### 5. Semaphore Pattern

Controls how many goroutines can access a shared resource simultaneously. Use semaphores to:

- Limit concurrent database connections
- Throttle API requests and prevent resource exhaustion under heavy load

There are many more patterns and variations, but these are some of the most common ones used in Go applications.

Check out [Go Concurrency Patterns](https://go.dev/talks/2012/concurrency.slide#1) for more details.

## Race conditions 🏁

Race conditions occur when multiple goroutines access shared data at the same time while at least one of them is modifying it.

The outcome of this can be unpredictable and difficult to debug.

We can use the `-race` flag while running or testing our Go code to detect race conditions.

```go
  go run -race main.go
  go test -race
```

This is helpful in identifying potential race conditions during development and testing phases. It's highly recommended and helps in writing safe concurrent code.

## Concurrency Implementation 🚧

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

// calculateSquare simulates a time-consuming calculation
func calculateSquare(num int) int {
	time.Sleep(time.Second)
	return num * num
}

// Without concurrency: execution is sequential
func sequentialCalculation(numbers []int) {
	fmt.Println("Sequential Calculation")
	start := time.Now()

	for _, num := range numbers {
		result := calculateSquare(num)
		fmt.Printf("%d² = %d\n", num, result)
	}

	fmt.Printf("Time taken: %v\n", time.Since(start))
}

// With concurrency: using goroutines and WaitGroup
func concurrentCalculation(numbers []int) {
	fmt.Println("Concurrent Calculation")
	start := time.Now()

	var wg sync.WaitGroup

	for _, num := range numbers {
		wg.Add(1)

		go func(n int) {
			defer wg.Done()
			result := calculateSquare(n)
			fmt.Printf("%d² = %d\n", n, result)
		}(num)
	}

	wg.Wait()
	fmt.Printf("Time taken: %v\n", time.Since(start))
}

// With channels: using goroutines and channels
func concurrentWithChannels(numbers []int) {
	fmt.Println("Concurrent with Channels")
	start := time.Now()

	results := make(chan string, len(numbers))

	for _, num := range numbers {
		go func(n int) {
			result := calculateSquare(n)
			results <- fmt.Sprintf("%d² = %d", n, result)
		}(num)
	}

	for range numbers {
		fmt.Println(<-results)
	}

	fmt.Printf("Time taken: %v\n", time.Since(start))
}

func main() {
	numbers := []int{1, 2, 3, 4, 5}

	sequentialCalculation(numbers) // Takes ~5 seconds
	fmt.Println("================")
	concurrentCalculation(numbers) // Takes ~1 second
	fmt.Println("================")
	concurrentWithChannels(numbers) // Takes ~1 second
}
```

- In the above example, we have three functions demonstrating different approaches to **calculate the square of numbers** in a slice.
- The `sequentialCalculation` function performs the calculations one after the other, taking around 5 seconds for 5 numbers.
- The `concurrentCalculation` function uses goroutines and a WaitGroup to perform the calculations concurrently, reducing the time to around 1 second.
- The `concurrentWithChannels` function uses goroutines and channels to achieve the same concurrent behavior, also taking around 1 second.
- Both concurrent approaches significantly improve performance by leveraging WaitGroups and channels to manage goroutines effectively.

> P.S: When running concurrently, the **output order** is **not guaranteed** because goroutines execute independently and may complete at different times.

## "Handle with care": Concurrency pitfalls ⚠️

- **Goroutine leaks:** When goroutines are not properly terminated and continue to run in the background, consuming resources. This can happen if they are blocked on a channel that is never written to or read from, or if they enter an infinite loop.
- **Deadlocks:** When using unbuffered channels, if two or more goroutines are waiting for each other to send or receive data, they can get stuck in a deadlock situation where none of them can proceed.
- **Race conditions:** When multiple goroutines access shared data simultaneously without ensuring proper synchronization, leading to unpredictable behavior.
- **Not waiting for goroutines to finish:** If the main goroutine exits before other goroutines have completed their work, those goroutines will be terminated prematurely.

To avoid these pitfalls, it's important to carefully design and test concurrent code, and to use synchronization mechanisms correctly.

> P.S Do not use concurrency unless it's really necessary!

## Conclusion 📝

- `Concurrency` is a powerful feature of Go that allows us to build responsive, fast and performant applications.
- `Goroutines` and `channels` are at the core of Go's concurrency model.
- More than understanding the concepts, it's important to know the pitfalls and best practices while working with concurrency.
- Using `concurrency patterns` can help in managing complexity and writing safe concurrent code.
