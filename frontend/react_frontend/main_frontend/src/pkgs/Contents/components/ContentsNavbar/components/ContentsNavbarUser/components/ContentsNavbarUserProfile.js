import baseProfileImg from "../assets/img/baseProfileImg.png"

import { useContentsNavbarUserGetDataHook } from "../hooks"

export const ContentNavbarUserProfile = ({ computerNumber, setComputerNumber }) => {
  const { profileImg, userDetail } = useContentsNavbarUserGetDataHook(computerNumber, setComputerNumber)

  return (
    <div className = "contentNavbarUserProfileContainer">
      
      {/* 프로필 이미지가 들어가는 장소 */}
      <div className = "contentNavbarUserProfileUserImgDiv">
        <img 
          src = { profileImg || baseProfileImg }
          className = "contentNavbarUserProfileUserImgPicture"
          alt = "유저이미지"
        />
      </div>

      {/* 유저의 정보가 담기는 장소 */}
      <div className = "contentNavbarUserProfileDetailDiv">
        <h3 className = "contentNavbarUserProfileDetailValue">
          { userDetail }
        </h3>
      </div>

    </div>
  )
}