# async
A 2-tier channel system. The first tier handles queuing tasks. The second controls the number workers that concurrently operate on the task queue. Code is based on the [post](http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/) on handling one million requests per minute with Go.
