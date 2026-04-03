# task api but in golang

### 3 april 2026

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

## stack

just **go**, that's it.

## find me

[portfolio](https://tgr-wjya.github.io) · [linkedin](https://linkedin.com/in/tegar-wijaya-kusuma-591a881b9) · [email](mailto:tgr.wjya.queue.top126@pm.me)

---

made with ◉‿◉