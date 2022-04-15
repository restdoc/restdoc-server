var app =
  new Vue(
    {
      el: "#app",
      delimiters:
        [
          "[[",
          "]]",
        ],
      data: {
        cur: 2,
        EntAdmin: 0,
        EntNormal: 1,
        Personal: 2,
        serverError: false,
        serverMessage: "",
        emailPattern: /^([A-Za-z0-9_\-\.\u4e00-\u9fa5])+\@([A-Za-z0-9_\-\.])+\.([A-Za-z]{2,8})$/,
        usernamePattern: /^([A-Za-z0-9])([A-Za-z0-9]{3,19})$/,
        personalUsername: "",
        personalUsernameRequiredError: false,
        personalUsernameFormatError: false,
        personalEmail: "",
        personalEmailRequiredError: false,
        personalEmailFormatError: false,
        personalPassword: "",
        passwordPattern: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[^]{8,30}$/,
        personalPasswordRequiredError: false,
        personalPasswordFormatError: false,
        personalCode: "",
        personalCodeRequiredError: false,
        entAdminUsername: "",
        entAdminUsernameRequiredError: false,
        entAdminUsernameFormatError: false,
        entAdminEmail: "",
        entAdminEmailRequiredError: false,
        entAdminEmailFormatError: false,
        entAdminPassword: "",
        entAdminPasswordRequiredError: false,
        entAdminPasswordFormatError: false,
        entAdminCode: "",
        entAdminCodeRequiredError: false,
        entNormalUsername: "",
        entNormalUsernameRequiredError: false,
        entNormalUsernameFormatError: false,
        entNormalEmail: "",
        entNormalEmailRequiredError: false,
        entNormalEmailFormatError: false,
        entNormalPassword: "",
        entNormalPasswordRequiredError: false,
        entNormalPasswordFormatError: false,
        entNormalCode: "",
        entNormalCodeRequiredError: false,
      },
      mounted: function(){
        this.initGeetestModule();
      },
      methods: {
            changeCur(cur) {
                this.cur = cur;
                this.serverError = false;
                this.serverMessage = "";
            },
          preprocessGeetest: function(el, type) {
                console.log('prepare geetest');
                var data = {};
                var _self = this;
                this.$http.get("/gt_preprocess?t=" + new Date().getTime(), data, { emulateJSON: true, credentials: true, }).then(
                    (response) => {
                        console.log(response);
                        initGeetest({
                                https: true,
                                gt: response.data.gt,
                                challenge: response.data.challenge,
                                new_captcha: response.data.new_captcha,
                                product: "embed", // 产品形式，包括：float，embed，popup。注意只对PC版验证码有效
                                offline: !response.data.success, // 表示用户后台检测极验服务器是否宕机，一般不需要关注
                                // 更多配置参数请参见：http://www.geetest.com/install/sections/idx-client-sdk.html#config
                            },
                            function(captchaObj) {
                                _self.handleEmbed(_self, captchaObj, type);
                            }
                        );
                        $(el).val(response.data.success);
                    }
                );
            },
            handleEmbed: function(_self, captchaObj, type) {
                var el = "";
                var waitEl = "";
                var captcha = {};
                console.log("cap");
                console.log(captchaObj);
                // 将验证码加到id为captcha的元素里，同时会有三个input的值：geetest_challenge, geetest_validate, geetest_seccode
                switch (type) {
                    case _self.EntAdmin:
                        el = "#embed-captcha-ent-admin";
                        waitEl = "#wait-ent-admin";
                        _self.ent_admin_captcha = captchaObj;
                        captcha = _self.ent_admin_captcha;
                        break;
                    case _self.EntNormal:
                        el = "#embed-captcha-ent-normal";
                        waitEl = "#wait-ent-normal";
                        _self.ent_normal_captcha = captchaObj;
                        captcha = _self.ent_normal_captcha;
                        break;
                    case this.Personal:
                        el = "#embed-captcha-personal";
                        waitEl = "#wait-personal";
                        _self.personal_captcha = captchaObj;
                        captcha = _self.personal_captcha;
                        break;
                    default:
                        return
                }
                captcha.appendTo(el);
                captcha.onReady(function() {
                    $(waitEl)[0].className = "hide";
                });

                // 更多接口参考：http://www.geetest.com/install/sections/idx-client-sdk.html
    
            },
            initGeetestModule: function() {
                //this.preprocessGeetest("#status-ent-admin", this.EntAdmin);
                //this.preprocessGeetest("#status-personal", this.Personal);
            },
            checkCaptcha() {
              return true;
              /*
                var validate = false;
                var el = "";

                switch (this.cur) {
                    case this.EntAdmin:
                        validate = this.ent_admin_captcha.getValidate();
                        el = "#notice-ent-admin";
                        break;
                    case this.EntNormal:
                        validate = this.ent_normal_captcha.getValidate();
                        el = "#notice-ent-normal";
                        break;
                    case this.Personal:
                        validate = this.personal_captcha.getValidate();
                        el = "#notice-personal";
                        break;
                    default:
                        return;
                }
                this.validateResult = validate;
                if (!validate) {
                    $(el)[0].className = "show";
                    setTimeout(function() { $(el)[0].className = "hide"; }, 2000);
                    return false;
                }
                return true;
                */
            },
            focusEmail: function() {
                this.entAdminEmailRequiredError = false;
                this.entAdminEmailFormatError = false;
            },
            checkEmail: function() {
                if (this.entAdminEmail == "") {
                    this.entAdminEmailRequiredError = true;
                    return false;
                } else {
                    this.entAdminEmailRequiredError = false;
                }
                if (!this.emailPattern.test(this.entAdminEmail)) {
                    this.entAdminEmailFormatError = true;
                    return false;
                }
                return true;
            },
            checkForm: function(needCheckCode) {
                var passed = this.checkUsername();
                if (!passed) {
                    console.log("user name");
                    return passed;
                }

                passed = this.checkPassword();
                if (!passed) {
                    console.log("password");
                    return passed;
                }

                passed = this.checkEmail();
                if (!passed) {
                    console.log("email");
                    return passed;
                }

                if (needCheckCode) {
                    passed = this.checkCode();
                    if (!passed) {
                        return passed;
                    }
                }
                return true;
            },
            checkCode: function () {
              if (this.entAdminCode == "") {
                  this.entAdminCodeRequiredError = true;
                  return false;
              } else {
                  this.entAdminCodeRequiredError = false;
              }
              return true;
            },
            getMailCode: function() {

                var result = this.checkForm(false);
                if (result != true) {
                    return;
                }

                var name = "";
                var email = "";
                var password = "";
              var url = "/forgotpassword/user";
                this.entAdminCodeRequiredError = false;
                name = this.entAdminUsername;
                email = this.entAdminEmail;
                password = this.entAdminPassword;
            
                result = this.checkCaptcha();
                if (result != true) {
                    return;
                }

                /*
                var geetest_challenge = this.validateResult.geetest_challenge;
                var geetest_validate = this.validateResult.geetest_validate;
                var geetest_seccode = this.validateResult.geetest_seccode;
                */

                var data = {
                    name: name,
                    company: this.company,
                    password: password,
                    email: email,
                    geetest_challenge: "",
                    geetest_validate: "",
                    geetest_seccode: "",
                };

                this.$http.post(url, data, { emulateJSON: true, credentials: true, }).then((response) => {
                        var err = response.data.error;
                        if (err && err != "") {
                            this.serverError = true;
                            if (err.includes("Error 1062: Duplicate entry")) {
                                this.serverMessage = "获取验证码失败";
                            } else {
                                this.serverMessage = "获取验证码失败 请稍候重试 " + err;
                            }
                        } else {
                            //
                        }
                    },
                    (response) => {
                        var err = response.data.error;
                        console.log(err);
                        this.serverError = true;
                        this.serverMessage = err;
                        //this.username = _username;
                        // error callback
                    }
                );
            },
            focusPassword: function() {
                        this.entAdminPasswordRequiredError = false;
                        this.entAdminPasswordFormatError = false;
            },
            checkPassword: function() {
                    if (this.entAdminPassword == "") {
                        this.entAdminPasswordRequiredError = true;
                        return false;
                    } else {
                        this.entAdminPasswordRequiredError = false;
                    }
                    if (!this.passwordPattern.test(this.entAdminPassword)) {
                        this.entAdminPasswordFormatError = true;
                        return false;
                    }
                    return true;
                    
            },
            focusCode: function() {
                        this.entAdminCodeRequiredError = false;
            },
            focusUsername: function() {
                        this.entAdminUsernameRequiredError = false;
                        this.entAdminUsernameFormatError = false;
                        
            },
             checkNameError: function(_self) {
                        if (_self.entAdminUsername == "") {
                            _self.entAdminUsernameRequiredError = true;
                            return false;
                        }
                        _self.entAdminUsernameRequiredError = false;
                        if (_self.usernamePattern.test(_self.entAdminUsername)) {
                            _self.entAdminUsernameFormatError = false;
                            return true;
                        }
                        _self.entAdminUsernameFormatError = true;
                        return false;

            },
            checkUsername: function() {
                var _self = this;
                var result = this.checkNameError(_self);
                return result;
            },
            handlerEnterpriseEmbed: function (captchaObj) {
              this.enterprise_captcha =
                captchaObj;

              // 将验证码加到id为captcha的元素里，同时会有三个input的值：geetest_challenge, geetest_validate, geetest_seccode
              this.enterprise_captcha.appendTo(
                "#embed-captcha-enterprise"
              );
              this.enterprise_captcha.onReady(
                function () {
                  $(
                    "#wait-enterprise"
                  )[0].className =
                    "hide";
                }
              );
              // 更多接口参考：http://www.geetest.com/install/sections/idx-client-sdk.html
            },
          handlerPersonalEmbed:
            function (
              captchaObj
            ) {
              this.personal_captcha =
                captchaObj;

              // 将验证码加到id为captcha的元素里，同时会有三个input的值：geetest_challenge, geetest_validate, geetest_seccode
              this.personal_captcha.appendTo(
                "#embed-captcha-personal"
              );
              this.personal_captcha.onReady(
                function () {
                  $(
                    "#wait-personal"
                  )[0].className =
                    "hide";
                }
              );
              // 更多接口参考：http://www.geetest.com/install/sections/idx-client-sdk.html
            },
          checkPersonalPassword:
            function () {
              if (
                this
                  .personalPassword ==
                ""
              ) {
                this.personalPasswordRequiredError = true;
              } else {
                this.personalPasswordRequiredError = false;
              }

              if (
                !this.passwordPattern.test(
                  this
                    .personalPassword
                )
              ) {
                this.personalPasswordFormatError = true;
                return;
              }
              this.checkPersonalRepeatedpassword();
            },
          checkEnterprisePassword:
            function () {
              if (
                this
                  .enterprisePassword ==
                ""
              ) {
                this.enterprisePasswordRequiredError = true;
                return false;
              } else {
                this.enterprisePasswordRequiredError = false;
              }

              if (
                !this.passwordPattern.test(
                  this
                    .enterprisePassword
                )
              ) {
                this.enterprisePasswordFormatError = true;
                return false;
              }
              return true;
            },
          resetPassword: function ( e) {
              e.preventDefault();

              var result = this.checkForm(true);

              if ( result != true) {
                return;
              }

              var data = { email: this.email, };

              var url = "/resetpassword/user";
              data = { "email": this.entAdminEmail, "password": this.entAdminPassword, "code": this.entAdminCode };
                  
              this.$http.post(url, data, { emulateJSON: true, credentials: true, }).then( ( response) => {
                    var code = response.data.code;

                    if (code == 0) {
                      this.serverMessage = "";
                      this.serverError = false;
                      window.location.href = "/login";
                    } else {
                      var message = response.data.message;
                      this.serverMessage = message;
                      this.serverError = true;
                    }
                  }, ( response) => {
                    console.log( response);
                    this.serverMessage = "服务器开小差儿了";
                    this.serverError = true;
                  }
                );
            },
          checkPersonalForm:
            function () {
              if (this.personalUsernameRequiredError || this.personalUsernameFormatError) {
                return false;
              }
              return true;
            },
          checkEnterpriseForm:
            function () {
              if (
                this
                  .emailRequiredError ||
                this
                  .emailFormatError
              ) {
                return false;
              }
              return true;
            },
      },
      
    }
  );
