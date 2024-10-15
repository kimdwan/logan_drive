import { useNavigate } from "react-router-dom"
import { useCallback } from "react"

export const useMainLoginClickSignUpBtnHook = () => {
  const navigate = useNavigate()
  
  const clickSignUpBtn = useCallback((event) => {

    if (event.target.className === "mainLoginFormSignUpBtn") {
      event.preventDefault()
      navigate("/signup/term/")
    }

  }, [ navigate ])

  return { clickSignUpBtn }
}