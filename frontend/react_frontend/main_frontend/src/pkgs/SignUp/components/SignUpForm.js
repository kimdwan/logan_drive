import { useSignUpFormHook, useSignUpGetTermAgree3Hook } from "../hooks"

export const SignUpForm = () => {

  // 선택 사항에 어떻게 체크했는지 확인하는 로직 
  const { termAgree3 } = useSignUpGetTermAgree3Hook()

  // 폼에서 데이터를 뽑아오는 함수
  const { register, handleSubmit, errors, onSubmit } = useSignUpFormHook()

  return (
    <div className = "signUpFormContainer">
      <form onSubmit = {handleSubmit(onSubmit)}>
        
        {/* 입력이 이루어지는 장소 */}
        <div className = "signUpFormInputDivBox">

          {/* 이메일이 입력되는 장소 */}
          <div className = "signUpFormInputEmailDivBox">

            {/* 글씨 */}
            <div className = "signUpFormInputEmailWordBox">
              <div className = "signUpFormInputEmailWordSmallBox">
                <h2 className = "signUpFormInputEmailWordText">이메일</h2>
              </div>
            </div>

            {/* 값을 입력하는 장소 */}
            <div className = "signUpFormInputEmailValueBox">

              {/* 값이 나오는 장소 */}
              <div className = "signUpFormInputEmailValueInput">
                <input 
                  className = "signUpFormInputEmailValues"
                  type = "text"
                  { ...register("email") }
                />
              </div>

              {/* 에러값 배출 */}
              <div className = "signUpFormInputEmailErrorBox">
                {
                  errors.email?.message && <p className = "signUpFormErrorMsg">{errors.email.message}</p>
                }
              </div>
            </div>
          </div>

          {/* 닉네임이 입력되는 장소 */}
          <div className = "signUpFormInputNicknameDivBox">
              
            {/* 글씨 */}
            <div className = "signUpFormInputNicknameWordBox">
              <div className = "signUpFormInputNicknameWordSmallBox">
                <h2 className = "signUpFormInputNicknameWordText">닉네임</h2>
              </div>
            </div>
            
            {/* 값을 입력하는 장소 */}
            <div className = "signUpFormInputNicknameValueBox">
              
              {/* 값이 나오는 장소 */}
              <div className = "signUpFormInputNicknameValueInput">
                <input 
                  className = "signUpFormInputNicknameValues"
                  type = "text"
                  { ...register("nickname") }
                />
              </div>

              {/* 에러값 배출 */}
              <div className = "signUpFormNicknameErrorBox">
                {
                  errors.nickname?.message && <p className = "signUpFormErrorMsg">{ errors.nickname.message }</p>
                }
              </div>
            </div>
                
          </div>

          {/* 비밀번호가 입력되는 장소 */}
          <div className = "signUpFormInputPasswordDivBox">
                
            {/* 글씨 */}
            <div className = "signUpFormInputPasswordWordBox">
              <div className = "signUpFormInputPasswordWordSmallBox">
                <h2 className = "signUpFormInputPasswordWordText">비밀번호</h2>
              </div>
            </div>

            {/* 값을 입력하는 장소 */}
            <div className = "signUpFormInputPasswordValueBox">
              
              {/* 값이 나오는 장소 */}
              <div className = "signUpFormInputPasswordValueInput">
                <input 
                  className = "signUpFormInputPasswordValues"
                  type = "password"
                  { ...register("password") }
                />
              </div>

              {/* 에러값 배출 */}
              <div className = "signUpFormInputPasswordErrorBox">
                {
                  errors.password?.message && <p className = "signUpFormErrorMsg">{ errors.password.message }</p>
                }
              </div>
            </div>

          </div>

          {/* 비밀번호 확인이 입력되는 장소 */}
          <div className = "signUpFormInputConfirmPasswordDivBox">
            
            {/* 글씨 */}
            <div className = "signUpFormInputConfirmPasswordWordBox">
              <div className = "signUpFormInputConfirmPasswordWordSmallBox">
                <h2 className = "signUpFormInputConfirmPasswordText">비밀번호 확인</h2>
              </div>
            </div>

            {/* 값을 입력하는 장소 */}
            <div className = "signUpFormInputConfirmPasswordValueBox">
              
              {/* 값이 나오는 장소 */}
              <div className = "signUpFormInputConfirmPasswordValueInput">
                <input 
                  className = "signUpFormInputConfirmPasswordValues"
                  type = "password"
                  { ...register("confirm_password") }
                />
              </div>

              {/* 에러값 배출 */}
              <div className = "sigUpFormInputConfirmPasswordErrorBox">
                {
                  errors.confirm_password?.message && <p className = "signUpFormErrorMsg">{ errors.confirm_password.message }</p>
                }
              </div>        
            </div>
          </div>

          {/* 회원가입 버튼이 있는 장소 */}
          <div className = "signUpFormSubmitBtnBox">
            <button
              className = "signUpFormSubmitBox"
              type = "submit"
            >
              회원가입
            </button>
          </div>

        </div>

      </form>
    </div>
  )
}