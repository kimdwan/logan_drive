import { useEffect, useState } from "react"
import { useLocation } from "react-router-dom"

export const useContentNavbarUserServiceReedemLocationHook = () => {
  const browser = useLocation()
  const [ urlTitle, setUrlTitle ] = useState("")

  useEffect(() => {
    const hostName = browser.pathname
    const hostNameList = hostName.split("/")
    setUrlTitle(hostNameList[hostNameList.length - 1])
  }, [browser])

  return { urlTitle }
}