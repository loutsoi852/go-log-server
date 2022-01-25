import { useEffect } from 'react'
import './App.css';

function App() {

  var interval = 0
  async function post() {
    try {
      await fetch('http://' + window.location.hostname + ':7777/send', {
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
      const req = await fetch('http://' + window.location.hostname + ':7777/read/6', {
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

  console.log('ws connect')
  const socket = new WebSocket("ws://localhost:7777/liveLogs");

  socket.onmessage = function (e) {
    console.log('ws msg')
    console.log(e.data)
  };


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
