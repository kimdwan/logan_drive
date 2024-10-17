import { useCallback } from "react"
import { useNavigate } from "react-router-dom"

export const useSignUpNavbarClickLogoImgHook = () => {
  const navigate = useNavigate()

  const clickLogoImg = useCallback((event) => {
    if (event.target.className === "signUpNavbarLogoImg") {
      navigate("/")
    }

  }, [ navigate ])

  return { clickLogoImg }
}