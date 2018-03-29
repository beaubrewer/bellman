# Bellman
*at your service*

Bellman is a Google Calendar driven doorbell sound manager. Doobell themes are defined in the event description using YAML. Bellman really needs a cool logo!

### Installation

Bellman requires
- [A MQTT Broker](https://www.digitalocean.com/community/tutorials/how-to-install-and-secure-the-mosquitto-mqtt-messaging-broker-on-ubuntu-16-04) to run. [CloudMQTT](https://www.cloudmqtt.com/plans.html) offers a free tier if you prefer not to install one.
- A Google Account and Calendar to schedule sound themes

```sh
$ ./bellman configure
```
Once you have configured Bellman run...

```sh
$ ./bellman server
```

### Development

Want to contribute to Bellman? Great!
Let's work together to make the only doorbell manager... the best doorbell manager!

### Todos

 - Dockerize
 - Add sound events that play at specified date/times
 - Improve configuration

License
----
MIT License

Copyright (c) 2018 Beau Brewer

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

**Free Software for WildWorks VIPs :)**