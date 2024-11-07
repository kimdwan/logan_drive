import { MainSupport } from "./MainSupport"
import { useMainSupportClickMainLogoHook } from "../hooks"

import mainLogo from "../assets/img/mainLogo.webp"

export const MainNavbar = () => {
  const { clickMainLogo } = useMainSupportClickMainLogoHook()

  return (
    <div className = "mainNavbarContainer">
      
      {/* 네브바에 들어가는 이미지와 로고 */}
      <div className = "mainNavbarDivBox">  
        <img className = "mainNavbarImg" 
          src = { mainLogo } 
          alt = "메인로고" 
          style = {{ cursor : "pointer" }}
          onClick = { clickMainLogo }
        />
        {/* 모달 */}
        <MainSupport/>

        <div className = "mainNavbarDiv" style = { { cursor : "pointer" } }>
          <h2 className = "mainNavbarDivH2Value" onClick = { clickMainLogo }>공유저장소_</h2>
          <h3 className = "mainNavbarDivH3Value" onClick = { clickMainLogo }>LOGAN</h3>
        </div>
      </div>

    </div>
  )
}