import { useEffect, useState } from "react"
import { useLocation } from "react-router-dom"

export const useContentNavbarTopRightGetUrlNameHook = () => {
    const location = useLocation()
    const [ urlName, setUrlName ] = useState("")
    useEffect(() => {
      const url_link = location.pathname
      const path_lists = url_link.split("/")
      setUrlName(path_lists[path_lists.length - 1])
  
    }, [ location ])

    return { urlName }
}