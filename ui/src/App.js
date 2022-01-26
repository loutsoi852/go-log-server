import { useEffect } from 'react'
import './App.css';

function App() {
  const ip = "10.0.1.197"
  //const ip = window.location.hostname

  var interval = 0

  async function post() {
    try {
      await fetch('http://' + ip + ':7777/send', {
        method: 'POST', // *GET, POST, PUT, DELETE, etc.
        headers: {},
        body: JSON.stringify({
          log: "hello world " + (new Date()).toUTCString()
        })
      })
      console.log("post success")
    } catch (e) {
      console.log('post failed')
    }
  }

  function postMulti() {
    clearInterval(interval)
    interval = setInterval(() => post(), 500)
  }

  function stopMulti() {
    clearInterval(interval)
  }

  async function read() {
    try {
      const req = await fetch('http://' + ip + ':7777/read/6', {
        method: 'GET', // *GET, POST, PUT, DELETE, etc.
        headers: {},
      })
      const res = await req.json()
      console.log('res', res)
      console.log("read success")
    } catch (e) {
      console.log('read failed')
    }
  }
  const socket = new WebSocket("ws://" + ip + ":7777/liveLogs");

  socket.onopen = function () {
    console.log('ws connected')
  };

  socket.onmessage = function (e) {
    const json = JSON.parse(e.data)
    console.log(json.time, JSON.parse(json.log))
  }


  useEffect(() => {
    return () => {
      socket.close()
    }
  }, [])

  return (
    <div className="App">
      <div>
        <button onClick={() => post()}>POST</button>
        <button onClick={() => read()}>READ</button>
        <button onClick={() => postMulti()}>POST Multi</button>
        <button onClick={() => stopMulti()}>Stop Multi</button>
      </div>
    </div>
  );
}

export default App;
