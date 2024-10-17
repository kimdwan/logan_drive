import subLogo from "../assets/img/subLogo.webp"

import { useSignUpNavbarClickLogoImgHook } from "../hooks"

export const SignUpNavbar = () => {
  const { clickLogoImg } = useSignUpNavbarClickLogoImgHook()

  return (
    <div className = "signUpNavbarContainer">
      <div className = "signUpNavbarLogoBox">
        <img className = "signUpNavbarLogoImg" src = {subLogo} alt = "서브로고" onClick = {clickLogoImg} style={{ cursor : "pointer" }} />
      </div>
    </div>
  )
}