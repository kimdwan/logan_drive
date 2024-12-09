import { useEffect, useState } from "react"
import { useLocation } from "react-router-dom"
import { ContentNavbarUserServiceChannelList } from "./components"

export const ContentNavbarUserService = () => {

  const browser = useLocation()
  const [ urlTitle, setUrlTitle ] = useState("")

  useEffect(() => {
    const hostName = browser.pathname
    const hostNameList = hostName.split("/")
    setUrlTitle(hostNameList[hostNameList.length - 1])
  }, [browser])

  return (
    <div className = "contentNavbarUserServiceContainer">
      {
        urlTitle === "channellist" ? <ContentNavbarUserServiceChannelList /> : <></>
      }
    </div>
  )
}