import "../assets/css/ContentsNavbarUserServiceChannelList.css"

import downArrow from "../assets/img/downArrow.png"

export const ContentNavbarUserServiceChannelList = () => {

  return (
    <div className = "contentNavbarUserServiceChannelListContainer">

      {/* 채널 서비스 목록 */}
      <div className = "contentNavbarUserServiceChannelListServiceListDiv">

        {/* 제목 */}
        <div className = "contentNavbarUserServiceChannelListServiceListMenuTitleBox">
          {/* 이름이 있음 */}
          <div className = "contentNavbarUserServiceChannelListServiceListMenuTitleValueBox">
            <h3 className = "contentNavbarUserServiceChannelListServiceListMenuTitleValue">
              서비스  
            </h3>
          </div>

          {/* 화살표 */}
          <div className = "contentNavbarUserServiceChannelListServiceListMenuArrowBox">
            <img 
              className = "contentNavbarUserServiceChannelListServiceListMenuArrowValue"
              src = {downArrow}
              alt = "화살표 사진"
            />
          </div>
        </div>

        {/* 리스트 */}
        <div className = "contentNavbarUserServiceChannelListListSortBox">
          {
            Array.from(["채널 리스트", "채널 생성"]).map((val, idx) => {
              return (
                <div className = {`contentNavbarUserServiceChannelListSortValue ${idx}`} key = {idx}>
                  <h4 className = "contentNavbarUserServiceChannelListSortWord">
                    {
                      val
                    }
                  </h4>
                </div>
              )
            })
          }
        </div>

      </div>

      {/* 채널 목록 (추후 추가) */}
      <div className = "contentNavbarUserServiceChannelListVarietyDiv">
        
        {/* 제목 */}
        <div className = "contentNavbarUserServiceChannelListVarietyTitleBox">
          {/* 이름이 있음 */}
          <div className = "contentNavbarUserServiceChannelListVarietyTitleValueBox">
            <h3 className = "contentNavbarUserServiceChannelListVarietyTitleValue">
              채널들
            </h3>
          </div>

          {/* 화살표 */}
          <div className = "contentNavbarUserServiceChannelListVarietyArrowBox">
            <img 
              className = "contentNavbarUserServiceChannelListVarietyArrowValue"
              src = {downArrow}
              alt = "화살표 사진"
            />
          </div>
        </div>

        {/* 리스트 */}
        <div className = "contentNavbarUserServiceChannelListUserSettingBox">
          추후 추가
        </div>

      </div>

    </div>
  )
}