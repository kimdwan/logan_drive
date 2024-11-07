import { useCallback, useEffect, useRef } from "react"

export const useMainSupportClickMainLogoHook = () => {
  const modalDiv = useRef(null)
  const modalDivBox = useRef(null)

  const clickMainLogo = useCallback((event) => {
    if (event.target.className === "mainNavbarImg" || event.target.className === "mainNavbarDivH2Value" || event.target.className === "mainNavbarDivH3Value") {
      modalDiv.current = document.querySelector("#mainSupportMainModar")
      modalDiv.current.showModal()
      modalDivBox.current = document.querySelector(".mainSupportContainer")
      modalDivBox.current.style.display = "block"
    } else if ( event.target.className === "mainSupportModarCloseBtn" ) {
      modalDiv.current = document.querySelector("#mainSupportMainModar")
      modalDiv.current.close()
      modalDivBox.current = document.querySelector(".mainSupportContainer")
      modalDivBox.current.style.display = "none" 
    } 
  }, [ ])

  useEffect((event) => {

  },[])

  return { clickMainLogo }
}