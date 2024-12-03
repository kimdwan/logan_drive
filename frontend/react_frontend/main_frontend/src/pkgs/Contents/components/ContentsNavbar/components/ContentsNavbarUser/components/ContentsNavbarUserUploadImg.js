import { useContentsNavbarUploadImgHook } from "../hooks"

export const ContentNavbarUserUploadImg = ({ computerNumber, setComputerNumber }) => {
  const { register, handleSubmit, onSubmit, clickUploadImgBtn, userImgData } = useContentsNavbarUploadImgHook(computerNumber, setComputerNumber)

  return (
    <div className = "contentNavbarUserUploadImgContainer">

      {/* 프로필 이미지를 업로드 하는 장소 */}
      <form encType = "multipart/form-data" onSubmit = { handleSubmit(onSubmit) }>

        {/* 실질적으로 프로필 이미지가 업로드 됨 */}
        <div className = "contentNavbarUserUploadImgConveyDiv">
          <input 
            id = "contentNavbarUserUploadImgConveyInterectBox"
            type = "file"
            {
              ...register("user_profile_img")
            }
            hidden
            onChange = {clickUploadImgBtn}
          />
          <label 
            className = "contentNavbarUserUploadImgConveyInterectSubstitute"
            htmlFor = "contentNavbarUserUploadImgConveyInterectBox"
            style = {{ cursor : "pointer" }}
          >
            <h2 className = "contentNavbarUserUploadImgConveyInterecSubstituteValue">이미지 업로드</h2>
          </label>
        </div>
        
        {/* 무슨 사진을 업로드 하는지 보이는 장소 */}
        <div className = "contentNavbarUserUploadImgPerformDiv">

          {/* 사진이 보이는 장소 */}
          <div className = "contentNavbarUserUploadImgInvestigateDiv">
            <h3 className = "contentNavbarUserUploadImgInvestigateValue">{ userImgData }</h3>
          </div>

          {/* 사진 (사진이 있어야만 활성화) 올리기 버튼 */}
          <div className = "contentNavbarUserUploadImgArgueDiv">
            <input 
              className = "contentNavbarUserUploadImgArgueInputBox"
              type = "submit"
              value = "업로드"
            />
          </div>

        </div>
      </form> 

    </div>
  )
}