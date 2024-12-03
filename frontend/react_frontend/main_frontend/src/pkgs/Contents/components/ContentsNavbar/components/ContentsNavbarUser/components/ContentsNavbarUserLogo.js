

export const ContentNavbarUserLogo = () => {
  return (
    <div className = "contentNavbarUserLogoContainer">
      
      {/* 실질적으로 값이 들어가는 장소 */}
      <div className = "contentNavbarUserLogoDiv">

        {/* 이름 */}
        <div className = "contentNavbarUserLogoTitleBox">
          <h2 className = "contentNavbarUserLogoTitleValue">공유 저장소</h2>
        </div>

        {/* 제작자 */}
        <div className = "contentNavbarUserLogoNameBox">
          <h3 className = "contentNavbarUserLogoNameValue">logan_kim</h3>
        </div>
      </div>

    </div>
  )
}