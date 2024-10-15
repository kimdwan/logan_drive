import { useSignUpCheckOneTwoValueHook, useSignUpTermAllCheckBoxClickHook } from "../hooks"

export const SignUpTerm = () => {

  // 전체 동의가 작동하게 해주는 훅
  const { changeAllCheckBox } = useSignUpTermAllCheckBoxClickHook()
  
  // 필수 1,2가 체크되게 해줌
  const { clickGoSignUpFormBtn } = useSignUpCheckOneTwoValueHook()

  return (
    <div className = "signUpTermContainer">
      
      {/* 이용약간이 적혀있는 박스 */}
      <div className = "signUpTermTopDivBox">
        <h2 className = "signUpTermTopH2Value">이용약간</h2>
      </div>

      {/* 전체 동의를 누를수 있는 박스 */}
      <div className = "signUpTermAllcheckBoxDiv">
        {/* 글씨 */}
        <div className = "signUpTermAllCheckBoxH2Div">
          <h2 className = "signUpTermAllCheckBoxH2Value">전체동의</h2>
        </div>
        {/* 체크 표시 */}
        <div className = "signUpTermAllCheckBoxInputDiv">
          <input 
            type = "checkbox"
            className = "signUpTermAllCheckBoxInputValue"
            onChange = {changeAllCheckBox}
          />
        </div>
      </div>

      {/* 하나씩 적혀 있는 체크박스 */}
      <div className = "signUpTermCheckBoxDiv">
        
        {/* 필수 1 */}
        <div className = "signUpTermCheckBoxTermOneDiv">

          {/* 제목과 체크박스 */}
          <div className = "signUpTermCheckBoxTermOneTitleDiv">
            {/* 글씨 */}
            <div className = "signUpTermCheckBoxTermOneH2Div">
              <h2 className = "signUpTermCheckBoxTermOneH2Value">
                필수1
              </h2>
            </div>
            {/* 체크 표시 */}
            <div className = "signUpTermCheckBoxTermOneInputDiv">
              <input 
                id = "signUpTermCheckBoxTermOneInputValue"
                className = "signUpTermCheckBoxTermInputValue"
                type = "checkbox"
              />
            </div>
          </div>

          {/* 내용 */}
          <div className = "signUpTermCheckBoxTermOneTextDiv">
            <div className = "signUpTermCheckBoxTermOneTextValue">
              <h5>
                필수 1에 대한 내용들....
              </h5>
            </div>
          </div>

        </div>

        {/* 필수 2 */}
        <div className = "signUpTermCheckBoxTermTwoDiv">

          {/* 제목과 체크박스 */}
          <div className = "signUpTermCheckBoxTermTwoTitleDiv">
            {/* 글씨 */}
            <div className = "signUpTermCheckBoxTermTwoH2Div">
              <h2 className = "signUpTermCheckBoxTermTwoH2Value">
                필수2
              </h2>
            </div>
            {/* 체크 표시 */}
            <div className = "signUpTermCheckBoxTermTwoInputDiv">
              <input 
                id = "signUpTermCheckBoxTermTwoInputValue"
                className = "signUpTermCheckBoxTermInputValue"
                type = "checkbox"
              />
            </div>  
          </div>

          {/* 내용 */}
          <div className = "signUpTermCheckBoxTermTwoTextDiv">
            <h5 className = "signUpTermCheckBoxTermTwoTextValue">
              필수 2에 대한 내용들.....
            </h5>
          </div>

        </div>

        {/* 선택 3 */}
        <div className = "signUpTermCheckBoxTermThreeDiv">
          
          {/* 제목과 체크박스 */}
          <div className = "signUpTermCheckBoxTermThreeTitleDiv">
            {/* 글씨 */}
            <div className = "signUpTermCheckBoxTermThreeH2Div">
              <h2 className = "signUpTermCheckBoxTermThreeH3Value">
                선택
              </h2>
            </div>
            {/* 체크 표시 */}
            <div className = "signUpTermCheckBoxTermThreeInputDiv">
              <input 
                id = "signUpTermCheckBoxTermThreeInputValue"
                className = "signUpTermCheckBoxTermInputValue"
                type = "checkbox"
              />
            </div>
          </div>

          {/* 내용 */}
          <div className = "signUpTermCheckBoxTermThreeTextDiv">
            <h5 className = "signUpTermCheckBoxTermThreeTextValue">
              선택에 대한 내용들.....
            </h5>
          </div>
        </div>

      </div>

      {/* 동의 버튼 */}
      <div className = "signUpTermSubmitBtnDiv">
        <button 
        type = "button"
        className = "signUpTermSubmitBtnValue"
        onClick = {clickGoSignUpFormBtn}
        >
          동의합니다.
        </button>
      </div>

    </div>
  )
}