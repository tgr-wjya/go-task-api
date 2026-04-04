# task api but in golang

## 4 april 2026

> yes, i'm rewriting my task-api project in golang

i've had enough with the **js/ts** slander, its time to actually see if the haters really have something to say  about their languages, especially **go**.

go thankfully is somewhat readable to me, until it isn't. this is just to relieve the itch of wanting to work with go.

either i'd quit midway or falling in love with this language. thankfully, i've had plenty of experience with elysia/hono and bun.

so, i'm not going in fully blind, wish me luck.

## live url

i really am not planning to deploy this at all, but i might change my mind.

## endpoints

| method | what it does | body |
| --- | --- | --- |
| `GET  /tasks/all` | return all tasks | — |
| `POST /tasks` | create a task | `title: string`, `status?: enum` |
| `GET  /tasks/:id` | return task by `id` | — |
| `PATCH  /tasks/:id` | update task | `title?: string`, `status?: enum` |
| `DELETE  /tasks/:id` | delete task | — |

## what i learned from go

- this is so cursed but go doesn't have a specific `export` keyword so if you want to export your `func`, `struct`, whatever.
- you literary have to capitalize it so it could be exported.

  - ```go
    // Now, you could export it to other package
    type Body struct {
      Greet string `json:"greet"`
    }

    // This won't work because it's unexported (private to the package)
    type body struct {
      Greet string `json:"greet"`
    }
    ```

- just remember these:
  - **capitalized**: exported (accessible to other packages + visible to encoding/json)
  **lowercase**: unexported (private to the package)

- this is already explained in the [docs](/docs/GO-REFS.md) but `w` and `r` are your write and read context equivalent to hono's `c`
- even for passing a basic json body, you need to initialize a `struct` for it because go doesn't support native `var` initialization with json.

- you could do this in JS/TS to initialize a JSON var.

  - ```ts
    let body = { greet: "Hello, World" }
    ```

- but in Go, you need to be more explicit.

  - ```go
    type Body struct {
      Greet string `json:"greet"`
    }

    body := Body{
      Greet: "Hello, World",
    }
    ```

- i mean, you could still pass a json body without writing a struct for it, using `map`. but with struct, its much more preferable.

  - ```go
    body := map[string]interface{}{
      "greet": "hello, world",
    }
    ```

- about `mux`, it matches incoming request paths
and dispatches them to the correct handler.
- specifically acting as your router.
  
  - ```go
    mux.HandleFunc("GET /tasks/all", tasks.GetAll)
    ```

- `:=` just mean initialize said variable with automatic type inference.
- it only lives inside a `func`.

  - ```go
    // These lives inside your func
    x := 10

    // And you're only allowed to use the normal var at package level.
    var x int = 10    
    ```

- to initialize port in go, you need to use `:` at the start of it like `:8080`

## stack

just **go**, that's it.

## find me

[portfolio](https://tgr-wjya.github.io) · [linkedin](https://linkedin.com/in/tegar-wijaya-kusuma-591a881b9) · [email](mailto:tgr.wjya.queue.top126@pm.me)

---

made with ◉‿◉
