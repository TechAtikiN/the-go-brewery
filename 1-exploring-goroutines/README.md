# Exploring Goroutines in Go 🏗️
Goroutines are Go's fundamental building blocks for concurrent programming, allowing multiple pieces of logic to execute simultaneously while providing elegant communication mechanisms between them. 

In programming, when we need to handle multiple operations at once, we typically talk about two concepts: concurrency and parallelism. Let’s explore what each means.

## What is Concurrency? 
- In a concurrent world, we're going to have **one stream of execution**, but we're going to have multiple tasks that are being executed at the same time. 
- For example, assume a chef who is making multiple dishes at the same time. The chef is the main process, and each dish is a `goroutine`. The chef can start preparing one dish, then while that dish is cooking, he can start cutting vegetables for another dish, and so on. This way, the chef is able to utilize his time more effectively by working on multiple tasks simultaneously.
- In Go, concurrency is achieved through **goroutines** and **channels**. We'll dive further into goroutines below.
- In concurrency, we can't rely on the order of execution of tasks, as they can be executed in any order, so we need to make sure they are synchronized.

## What is Parallelism? 
- In a parallel world, we have **multiple streams of execution**, and each stream can run on a different CPU core. This allows us to utilize all available resources more effectively.
- For example, assume multiple chefs working in a kitchen, each preparing a different dish at the same time. Each chef is a separate process, and they can work independently of each other. This way, the kitchen can produce multiple dishes simultaneously, utilizing all available resources effectively.
- In Go, we achieve parallelism using **gomaxprocs**. That is basically going to define the number of CPU cores that we want to use for the program. By default, Go will use all available CPU cores, but we can limit it to a specific number using the runtime package.

## Difference between Concurrency and Parallelism? 
- Concurrency is about **dealing with lots of things at once** by managing and switching between tasks, making progress on them without necessarily running them simultaneously.
- Parallelism is about **doing lots of things at once** by running tasks simultaneously on multiple processors or cores.
- Go is designed for concurrency and can leverage parallelism when available hardware resources allows it.

## What is a `goroutine`? 🔀
- A `goroutine` is a **lightweight thread** managed by the Go runtime. It allows us to run functions concurrently, meaning that multiple functions can run at the same time without blocking each other.
- Goroutines are created using the `go` keyword followed by a function call. For example:
	```go
	go runThisFunction()
	```
- This will create a new `goroutine` that runs `runThisFunction()` concurrently with the main function.
- Goroutines are lightweight because they are **managed by the Go runtime**, which means that we don't have to worry about managing them ourselves. 


## How does a `goroutine` work? 📊
- Since goroutines are spinning out separate processes, that means they are performing multiple tasks at the same time, or they are running concurrently.
	```mermaid
	gantt
		title Goroutines Running Concurrently with Main
		dateFormat X
		axisFormat %s
		
		section Main
		Main Function    :active, main, 0, 10
		
		section Goroutines
		Goroutine 1      :g1, 2, 8
		Goroutine 2      :g2, 3, 9
		Goroutine 3      :g3, 1, 6
	```

- When we come across a `goroutine` in the main function, it's going to spin out a separate process from the main process, everything else will continue to run in the main process. When we encounter another `goroutine`, it will spin out another process, and so on. But when the **main process terminates, all the goroutines will also end**. 
- So, if we want to keep the main process running until all goroutines have finished executing, we need to have a blocking mechanism to wait for all the goroutines to finish.
- We'll see soon about how we can achieve this blocking mechanism using **`WaitGroup`**.


## What is `WaitGroup`? ⏳
- **`WaitGroup`** is a way for us to **block the main process* until all goroutines have finished executing. 
- Think of it like an internal counter that keeps track of the number of goroutines that are currently running. When the `WaitGroup` counter reaches zero, it means that all goroutines have finished executing. 
- `WaitGroup` is part of the sync package in Go, and it provides a simple way to **synchronize goroutines**.

### Methods of `WaitGroup`: ⚙️
- **Add(n int)**: increments the `WaitGroup` counter by n.
- **Done()**: decrements the `WaitGroup` counter by 1.
- **Wait()**: blocks until the `WaitGroup` counter is zero.

<br/>

> What is defer?
> - Defer is a keyword in Go that allows us to **delay the execution of a function until the surrounding function returns**.
> - It is commonly used with `WaitGroup` to ensure that the _Done()_ method is called when the `goroutine` finishes executing, even if an error occurs. This way, we can ensure that the `WaitGroup` counter is decremented correctly and the main process can exit safely.


### How does `WaitGroup` work? 
- When we create a `WaitGroup`, it starts with a counter of zero. 
	```go
	wg := sync.WaitGroup{}
	```

- When we add a `goroutine`, we call the Add() method to increment the counter by 1. 
	```go
	wg.Add(1)
	```
- When the `goroutine` finishes executing, it calls the Done() method to decrement the counter by 1. 
	```go
	wg.Done()
	```
- The main process can call the Wait() method to block until the counter reaches zero, indicating that all goroutines have finished executing. 
	```go
	wg.Wait()
	```

## Example of using goroutines and `WaitGroup` 
```go
func main() {
	var wg sync.WaitGroup 

	wg.Add(1) 
	go func() {
		defer wg.Done() 
		fmt.Println("Executing 1st goroutine")
	}() 

	wg.Add(1) 
	go func() {
		defer wg.Done()
		fmt.Println("Executing 2nd goroutine")
	}()

	wg.Wait() 
	fmt.Println("All goroutines finished executing")
}
```

### Explanation of the code: 📝
1. Create a `WaitGroup` that will be used to wait for goroutines to finish
2. Increment the `WaitGroup` counter by 1 for the first `goroutine`
3. Decrement the `WaitGroup` counter by 1 when this `goroutine` finishes
4. This is an anonymous function, so we need to invoke it immediately
5. Increment the `WaitGroup` counter by 1 for the second `goroutine`
6. Wait for all goroutines to finish executing

### Order of execution: ⏰
- The order of execution of goroutines is not guaranteed, as they can run concurrently and may finish at different times.
- So, in this example as well, the output may vary each time we run the program.

Output X:
```
Executing 1st goroutine
Executing 2nd goroutine
All goroutines finished executing
```

Output Y:
```
Executing 2nd goroutine
Executing 1st goroutine
All goroutines finished executing
```

## TL;DR: 🤝	
- Goroutines are lightweight "threads" managed by the Go runtime and they are created using the `go` keyword.
- Concurrency is about managing multiple tasks and switching between them, while parallelism is about running multiple tasks simultaneously on different CPU cores.
- Goroutines can communicate with each other using channels.
- `WaitGroup` is a way to synchronize goroutines and ensure that they finish executing before the main process exits.
