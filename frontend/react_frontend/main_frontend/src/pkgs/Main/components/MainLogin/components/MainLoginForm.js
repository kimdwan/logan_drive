import { useMainLoginFormHook } from "../hooks"

export const MainLoginForm = ({ setComputerNumber }) => {
  const { register, handleSubmit, errors, onSubmit } = useMainLoginFormHook(setComputerNumber)

  return (
    <div className = "mainLoginFormContainer">
      <form onSubmit = { handleSubmit(onSubmit) }>
        
        {/* 인풋 값이 들어가는 장소 */}
        <div className = "mainLoginFormDivBox">
          {/* 이메일 인풋 값이 들어가는 장소 */}
          <div className = "mainLoginFormEmailDivBox">

            {/* 유저가 입력을 하는 장소 */}
            <div className = "mainLoginFormEmailInputBox"> 
              {/* 앞에 Email 입력이라고 알려주는 함수 */}
              <div className = "mainLoginFormEmailTagBox">
                <h2>Email</h2>
              </div>
              
              {/* 실질적으로 입력하는 장소 */}
              <div className = "mainLoginFormEmailInput">
                <input 
                  id = "mainLoginFormEmailInputValue"
                  type = "text"
                  {
                    ...register("email")
                  }
                />
                <label
                  className = "mainLoginFormEmailInputLabel" 
                  htmlFor = "mainLoginFormEmailInputValue"
                />
              </div>
            </div>

            {/* 에러 값이 나오는 장소 */}
            <div className = "mainLoginFormEmailErrorBox">
              {
                errors.email?.message && <p className = "mainLoginFormErrorMsg">{errors.email.message}</p>
              }
            </div>
          </div>

          {/* 비밀번호 인풋 값이 들어가는 장소 */}
          <div className = "mainLoginFormPasswordDivBox">

            {/* 유저가 입력을 하는 장소 */}
            <div className = "mainLoginFormPasswordInputBox">
              {/* 앞에 password 입력이라고 알려주는 함수 */}
              <div className = "mainLoginFormPasswordTagBox">
                <h2>Password</h2>
              </div>

              {/* 실질적으로 입력하는 함수 */}
              <div className = "mainLoginFormPasswordInput">
                <input 
                  id = "mainLoginFormPasswordInputValue"
                  type = "password"
                  {
                    ...register("password")
                  }
                />
                <label 
                  className = "mainLoginFormPasswordLabel"
                  htmlFor = "mainLoginFormPasswordInputValue"
                />
              </div>
            </div>

            {/* 에러값이 나오는 장소 */}
            <div className = "mainLoginFormPasswordErrorBox">
              {
                errors.password?.message && <p className = "mainLoginFormErrorMsg">{ errors.password.message }</p>
              }
            </div>
          </div>
        </div>

        {/* 로그인 값이 전달되거나 회원가입 창으로 이동하게 해주는 로고 */}
        <div className = "mainLoginFormBtnBox">

          {/* 로그인이 이루어 지는 장소 */}
          <div className = "mainLoginFormLoginBtnBox">
            <input 
              className = "mainLoginFormLoginBtn"
              type = "submit"
              value = "로그인"
            />
          </div>

          {/* 회원가입으로 이동시켜주는 컴퍼넌트 */}
          <div className = "mainLoginFormSignUpBtnBox">
            <button className = "mainLoginFormSignUpBtn">
              회원가입
            </button>
          </div>

        </div>

      </form>
    </div>
  )
}