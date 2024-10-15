import mainLogo from "../assets/img/mainLogo.webp"

export const MainNavbar = () => {
  return (
    <div className = "mainNavbarContainer">
      
      {/* 네브바에 들어가는 이미지와 로고 */}
      <div className = "mainNavbarDivBox">  
        <img className = "mainNavbarImg" src = { mainLogo } alt = "메인로고" />
        <div className = "mainNavbarDiv">
          <h2 className = " mainNavbarDivH2Value">공유저장소_</h2>
          <h3 className = "mainNavbarDivH3Value">LOGAN</h3>
        </div>
      </div>

    </div>
  )
}