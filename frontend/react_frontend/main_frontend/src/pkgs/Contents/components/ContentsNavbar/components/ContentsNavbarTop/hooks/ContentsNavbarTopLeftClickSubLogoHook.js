import { useCallback } from "react"
import { useNavigate } from "react-router-dom"

export const useContentNavbarTopLeftClickSubLogoHook = () => {
  const navigate = useNavigate()
  const clickHomeBtn = useCallback((event) => {

    if (event.target.className === "contentNavbarTopLeftImgValue") {
      navigate("/main/channellist")
    }
  

  }, [ navigate ])

  return { clickHomeBtn }
}