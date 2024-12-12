import { useEffect, useState } from "react"
import { useLocation } from "react-router-dom"

export const useContentNavbarTopMidGetTitleNameHook = () => {
  const location = useLocation()
  const [ titleName, setTitleName ] = useState("")

  useEffect(() => {

    const path_name = location.pathname
    const path_lists = path_name.split("/")
    const path = path_lists[path_lists.length - 1]
    if (path === "channellist") {
      setTitleName("채널 리스트")
    } else {
      setTitleName("")
    }

  }, [ location ])
  

  return { titleName }
}