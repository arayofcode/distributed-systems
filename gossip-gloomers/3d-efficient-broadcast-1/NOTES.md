# Notes and Challenges

Updated code from 3c. The code seems to be working but isn't a solution. I'm seeing this error:
```
WARN [2024-12-07 12:51:12,778] n15 stdout - maelstrom.process Error!
java.lang.AssertionError: Assert failed: Invalid dest for message #maelstrom.net.message.Message{:id 168940, :src "n15", :dest "n10", :body {:message 874, :msg_id 5356, :type "broadcast"}}
(get queues (:dest m))
	at maelstrom.net$validate_msg.invokeStatic(net.clj:173)
	at maelstrom.net$validate_msg.invoke(net.clj:165)
	at maelstrom.net$send_BANG_.invokeStatic(net.clj:200)
	at maelstrom.net$send_BANG_.invoke(net.clj:188)
	at maelstrom.process$stdout_thread$fn__15878$fn__15879$fn__15881.invoke(process.clj:147)
	at maelstrom.process$stdout_thread$fn__15878$fn__15879.invoke(process.clj:146)
	at maelstrom.process$stdout_thread$fn__15878.invoke(process.clj:140)
	at clojure.core$binding_conveyor_fn$fn__5823.invoke(core.clj:2047)
	at clojure.lang.AFn.call(AFn.java:18)
	at java.base/java.util.concurrent.FutureTask.run(FutureTask.java:317)
	at java.base/java.util.concurrent.ThreadPoolExecutor.runWorker(ThreadPoolExecutor.java:1144)
	at java.base/java.util.concurrent.ThreadPoolExecutor$Worker.run(ThreadPoolExecutor.java:642)
	at java.base/java.lang.Thread.run(Thread.java:1575)
WARN [2024-12-07 12:51:12,780] n16 stdout - maelstrom.process Error!
java.lang.AssertionError: Assert failed: Invalid dest for message #maelstrom.net.message.Message{:id 168952, :src "n16", :dest "n11", :body {:message 452, :msg_id 7470, :type "broadcast"}}
(get queues (:dest m))
	at maelstrom.net$validate_msg.invokeStatic(net.clj:173)
	at maelstrom.net$validate_msg.invoke(net.clj:165)
	at maelstrom.net$send_BANG_.invokeStatic(net.clj:200)
	at maelstrom.net$send_BANG_.invoke(net.clj:188)
	at maelstrom.process$stdout_thread$fn__15878$fn__15879$fn__15881.invoke(process.clj:147)
	at maelstrom.process$stdout_thread$fn__15878$fn__15879.invoke(process.clj:146)
	at maelstrom.process$stdout_thread$fn__15878.invoke(process.clj:140)
	at clojure.core$binding_conveyor_fn$fn__5823.invoke(core.clj:2047)
	at clojure.lang.AFn.call(AFn.java:18)
	at java.base/java.util.concurrent.FutureTask.run(FutureTask.java:317)
	at java.base/java.util.concurrent.ThreadPoolExecutor.runWorker(ThreadPoolExecutor.java:1144)
	at java.base/java.util.concurrent.ThreadPoolExecutor$Worker.run(ThreadPoolExecutor.java:642)
	at java.base/java.lang.Thread.run(Thread.java:1575)
INFO [2024-12-07 12:51:12,783] jepsen node n24 - maelstrom.db Tearing down n24
WARN [2024-12-07 12:51:12,787] n24 stdout - maelstrom.process Error!
java.lang.AssertionError: Assert failed: Invalid dest for message #maelstrom.net.message.Message{:id 168984, :src "n24", :dest "n19", :body {:message 457, :msg_id 2780, :type "broadcast"}}
(get queues (:dest m))
	at maelstrom.net$validate_msg.invokeStatic(net.clj:173)
	at maelstrom.net$validate_msg.invoke(net.clj:165)
	at maelstrom.net$send_BANG_.invokeStatic(net.clj:200)
	at maelstrom.net$send_BANG_.invoke(net.clj:188)
	at maelstrom.process$stdout_thread$fn__15878$fn__15879$fn__15881.invoke(process.clj:147)
	at maelstrom.process$stdout_thread$fn__15878$fn__15879.invoke(process.clj:146)
	at maelstrom.process$stdout_thread$fn__15878.invoke(process.clj:140)
	at clojure.core$binding_conveyor_fn$fn__5823.invoke(core.clj:2047)
	at clojure.lang.AFn.call(AFn.java:18)
	at java.base/java.util.concurrent.FutureTask.run(FutureTask.java:317)
	at java.base/java.util.concurrent.ThreadPoolExecutor.runWorker(ThreadPoolExecutor.java:1144)
	at java.base/java.util.concurrent.ThreadPoolExecutor$Worker.run(ThreadPoolExecutor.java:642)
	at java.base/java.lang.Thread.run(Thread.java:1575)
```

I'm yet to figure out why this is failing right now. Most likely, the receiver node has been deleted by the time sender comes out of sleep and sends another message. 

## Issues
- (Potentially) Node deleted by the time sender comes out of sleep
- Requirements not met:
  - Messages-per-operation is below 30
  - Median latency is below 400ms
  - Maximum latency is below 600ms
- Potential duplicate messages being sent.

Potential solutions:
- Reduce time for sleep
- Maintain set of msg_id. If duplicate, don't broadcast.