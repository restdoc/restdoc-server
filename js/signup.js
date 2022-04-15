var app =
    new Vue({
        el: "#app",
        delimiters: [
            "[[",
            "]]",
        ],
        data: {
            entAdminSendText: "获取验证码",
            entNotSendmail: true,
            personalNotSendmail: true,
            cur: 0,
            emailNotice: "",
            EntAdmin: 0,
            EntNormal: 1,
            Personal: 2,
            personalAgree: false,
            enterpriseAgree: false,
            serverError: false,
            serverMessage: "",
            emailPattern: /^([A-Za-z0-9_\-\.\u4e00-\u9fa5])+\@([A-Za-z0-9_\-\.])+\.([A-Za-z]{2,8})$/,
            company: "",
            usernamePattern: /^([A-Za-z0-9])([_A-Za-z0-9]{4,19})$/,
            passwordPattern: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[^]{8,30}$/,
            entAdminUsername: "",
            entAdminUsernameRequiredError: false,
            entAdminUsernameFormatError: false,
            entAdminPassword: "",
            entAdminPasswordRequiredError: false,
            entAdminPasswordFormatError: false,
            entAdminEmail: "",
            entAdminEmailRequiredError: false,
            entAdminEmailFormatError: false,
            entAdminCode: "",
            entAdminCodeRequiredError: false,
            entAdminPlaceholder: "请输入邮箱验证码",
            entAdminGetMailDisabled: false,
            ent_admin_captcha: {},
            entAdminAgree: false,

            validateResult: {},
        },
        mounted: function() {
            this.initGeetestModule();
        },
        methods: {
            changeCur(cur) {
                this.cur = cur;
                this.serverError = false;
                this.serverMessage = "";
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

            countDown(time) {
                let that = this;

                time--;
                let timer = setTimeout(function () {
                    switch (that.cur) {
                    case that.EntAdmin:
                        that.entAdminSendText = time + "秒";
                        break;
                    case that.Personal:
                        that.personalSendText = time + "秒";
                        break;
                    default:
                    }

                    that.countDown(time);
                }, 1000);

                if (time == 0) {
                    that.entAdminGetMailDisabled = false;
                    that.entAdminSendText = "获取验证码";
                    that.entAdminPlaceholder = "请输入邮箱验证码";

                    clearInterval(timer);
                    return
                }
                
            },
            preprocessGeetest: function(el, type) {
                console.log('prepare geetest');
                var data = {};
                var _self = this;
                this.$http.get("/gt_preprocess?t=" + new Date().getTime(), data, { emulateJSON: true, credentials: true, }).then(
                    (response) => {
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
                el = "#embed-captcha-ent-admin";
                waitEl = "#wait-ent-admin";
                _self.ent_admin_captcha = captchaObj;
                captcha = _self.ent_admin_captcha;
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
            signup: function() {

                var result = false;
                var url = "/signup";
                var data = {};

                result = this.checkForm(true);
                data = {
                    name: this.entAdminUsername,
                    company: this.company,
                    password: this.entAdminPassword,
                    email: this.entAdminEmail,
                    code: this.entAdminCode,
                    admin: "true",
                };
                if (result != true) {
                    return;
                }


                this.$http.post(url, data, { emulateJSON: true, credentials: true, }).then((response) => {
                        this.entAdminCallback(response);
                    },
                    (response) => {
                        var err = response.data.error;
                        this.serverError = true;
                        this.serverMessage = err;
                        //this.username = _username;
                        // error callback
                    }
                );
            },
            entAdminCallback: function(response) {
                var err = response.data.error;
                if (err && err != "") {
                    this.serverError = true;
                    if (err.includes("Error 1062: Duplicate entry")) {
                        this.serverMessage = "帐号已经注册，请直接登录。";
                    } else {
                        this.serverMessage = "注册失败，请稍候重试 " + err;
                    }
                } else {
                    window.location.href = "/team";
                }
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

                result = this.checkCaptcha();
                if (result != true) {
                    return;
                }


                let time = 60;
                this.entAdminSendText = time + "秒";
                this.entAdminGetMailDisabled = true;
                this.countDown(time);

                var geetest_challenge = this.validateResult.geetest_challenge;
                var geetest_validate = this.validateResult.geetest_validate;
                var geetest_seccode = this.validateResult.geetest_seccode;

                var email = "";
                var password = "";
                var name = "";

                var data = {
                    name: name,
                    company: this.company,
                    password: password,
                    email: email,
                    geetest_challenge: geetest_challenge,
                    geetest_validate: geetest_validate,
                    geetest_seccode: geetest_seccode,
                };

                this.entAdminCodeRequiredError = false;
                email = this.entAdminEmail;
                password = this.entAdminPassword;
                data = {
                    name: name,
                    company: this.company,
                    password: password,
                    email: email,
                    geetest_challenge: geetest_challenge,
                    geetest_validate: geetest_validate,
                    geetest_seccode: geetest_seccode,
                };

              
                this.$http.post("/getmailcode", data, { emulateJSON: true, credentials: true, }).then((response) => {
                        var err = response.data.error;
                        err = "";
                        if (err && err != "") {
                            this.serverError = true;
                            if (err.includes("Error 1062: Duplicate entry")) {
                                this.serverMessage = "帐号已经注册，请直接登录。";
                            } else {
                                this.serverMessage = "注册失败，请稍候重试 " + err;
                            }
                        } else {
                            this.entNotSendmail = false;
                            $('#ent-admin-code').attr('placeholder', "邮件已发送")
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
            focusCode: function() {
                this.entAdminCodeRequiredError = false;
            },
        },
    });