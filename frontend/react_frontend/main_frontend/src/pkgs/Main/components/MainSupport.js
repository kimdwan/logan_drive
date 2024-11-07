import { useMainSupportClickMainLogoHook } from "../hooks"

import mainLogo from "../assets/img/mainLogo.webp"

export const MainSupport = () => {

  const { clickMainLogo } = useMainSupportClickMainLogoHook() 

  return (
    <div className = "mainSupportContainer" style = { { display : "none" } }>
      <dialog id = "mainSupportMainModar">
        
        {/* 본격적으로 보이는 화면 */}
        <div className = "mainSupportMainModarMainDiv">
          
          {/* 메인 로고 */}
          <div className = "mainSupportMainModarMainLogoDiv">
            <img 
              className = "mainSupportMainModarMainLogoImg"
              src = {mainLogo}
              alt = "메인로고"
            />
          </div>

          {/* 환영 메세지 */}
          <div className = "mainSupportMainModarMainWelcomeMsgDiv">
            <h2 className = "mainSupportMainModarMainWelcomeMsgValue">
              환영합니다.
            </h2>
          </div>

          {/* 나의 정보 */}
          <div className = "mainSupportMainModarMainDetail">
            
            {/* 이메일 */}
            <div className = "mainSupportMainModarMainEmailDiv">
              {/* 제목 */}
              <div className = "mainSupportMainModarMainEmailWordBox">
                <h3 className = "mainSupportMainModarMainEmailWordValue">
                  이메일
                </h3>
              </div>
              {/* 내용 */}
              <div className = "mainSupportMainModarMainEmailDetailBox">
                <h4 className = "mainSupportMainModarMainEmailDetailValue">
                  naxtto@naver.com / dongwan123456789@gmail.com 
                </h4> 
              </div>
            </div>

            {/* 은행 */}
            <div className = "mainSupportMainModarBankAccountDiv">
              {/* 제목 */}
              <div className = "mainSupportMainModarBankAccountWordBox">
                <h3 className = "mainSupportMainModarBankAccounrWordValue"> 
                  은행
                </h3>
              </div>
              {/* 내용 */}
              <div className = "mainSupportMainModarMainBankAccountDetailBox">
                <h4 className = "mainSupportMainModarMainBankAccountDetailValue">
                  307002-04-323703 국민은행
                </h4>
              </div>
            </div>

            {/* 후원시 좋은점 */}
            <div className = "mainSupportMainModarProfitDiv">
              {/* 제목 */}
              <div className = "mainSupportMainModarProfitWordBox">
                <h3 className = "mainSupportMainModarProfitWordValue">
                  후원시 혜택
                </h3>
              </div>
              {/* 내용 */}
              <div className = "mainSupportMainModarProfitDetailBox">
                <textarea 
                  readOnly
                  value = {`
                    1. 후원자 칭호 획득가능 (금액에 따라 요구 가능)
                    2. 추후 다른 웹사이트에서도 배네핏 적용 가능
                    3. 새로운 웹사이트를 만들때 의견 제시 가능
                    4. 오픈 채팅방 입장 가능 
                    # 주의사항: 후원 후 연락이 가능한 이메일을 알려주세요.
                    `}
                />
              </div>
            </div>

            {/* 닫기 버튼 */}
            <div className = "mainSupportModarCloseBtnDiv">
              <button className = "mainSupportModarCloseBtn" 
                onClick = {clickMainLogo}
              >
                닫기
              </button>
            </div>
          </div>
        </div>
      </dialog>
    </div>
  )
}