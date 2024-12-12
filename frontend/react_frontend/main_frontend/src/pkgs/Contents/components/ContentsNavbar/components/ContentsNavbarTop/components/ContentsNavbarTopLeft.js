import subMainLogo from "../assets/img/subMainLogo.webp"

import { useContentNavbarTopLeftClickSubLogoHook } from "../hooks"

export const ContentNavbarTopLeft = () => {
  const { clickHomeBtn } = useContentNavbarTopLeftClickSubLogoHook()

  return (
    <div className = "contentNavbarTopLeftContainer">
      
      {/* 로고가 있는 장소 */}
      <div className = "contentNavbarTopLeftImgDiv">
        {/* 로고 */}
        <img 
          onClick = { clickHomeBtn }
          className = "contentNavbarTopLeftImgValue"
          src = { subMainLogo }
          alt = "로고"
          style = {{ cursor : "pointer" }}
        />
      </div>

    </div>
  )
}