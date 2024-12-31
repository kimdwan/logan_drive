import { useEffect, useState, useRef } from "react"

export const useContentNavbarFriendGetStatusHook = (computerNumber) => {
  
    const [friendStatus, setFriendStatus] = useState([])
    const connectWs = useRef(null)
  
    useEffect(() => {
      const go_ws_url = process.env.REACT_APP_GO_WS_URL
  
      connectWs.current = new WebSocket(`${go_ws_url}/ws/user/status`)
      const ws = connectWs.current
      
      if (computerNumber) {
    
        ws.onopen = () => {
          console.log("연결 되었습니다.")
          ws.send(JSON.stringify({ Computer_number : computerNumber }))
        }
    
        ws.onmessage = (event) => {
    
          const messageData = JSON.parse(event.data)
          if (messageData) {
            setFriendStatus(messageData)
          }
    
        }
    
        ws.onclose = () => {
          console.log("닫혔습니다.")
        }
    
        ws.onerror = (event) => {
          alert("에러가 발생했습니다.")
          console.log(event)
        }
      }
  
      return () => {
  
        if ( ws.readyState === WebSocket.OPEN || ws.readyState === WebSocket.CONNECTING ) {
          ws.close()
        }
  
      }
  
    },[ computerNumber, setFriendStatus ])

    return { friendStatus }

}