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
        passwordPattern: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[^]{8,30}$/,
        entAdminUsername: "",
        entAdminUsernameRequiredError: false,
        entAdminUsernameFormatError: false,
        entAdminPassword: "",
        entAdminPasswordRequiredError: false,
        entAdminPasswordFormatError: false,
        personalUsername: "",
        personalPassword: "",
        personalUsernameRequiredError: false,
        personalUsernameFormatError: false,
        personalPasswordRequiredError: false,
        personalPasswordFormatError: false,
        admin: false,
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
            focusPassword: function() {
                this.entAdminPasswordRequiredError = false;
            },
            checkNameError: function(_self) {
                if (_self.entAdminUsername == "") {
                    _self.entAdminUsernameRequiredError = true;
                    return false;
                }
                _self.entAdminUsernameRequiredError = false;
                if (_self.emailPattern.test(_self.entAdminUsername)) {
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

                return true;
            },
            login: function () {
              var result = this.checkForm();
              if ( result != true) {
                return;
              }

              var data = {};
              var isAdmin = false;


              isAdmin = true;
              data = { email: this.entAdminUsername, password: this.entAdminPassword, admin: "true" };

              this.$http.post("/login", data, { emulateJSON: true, credentials: true, }).then( ( response) => {
                var code = response.data.code;
                console.log(response.data);
                if (code == 0) {
                  if (isAdmin) {
                    window.location.href = "/team";
                    this.serverMessage = "";
                    this.serverError = false;
                  } else {
                    window.location.href = "/restdoc";
                    this.serverMessage = "";
                    this.serverError = false;
                  }
                } else {
                    console.log("error");
                    var message = response.data.message;
                    console.log()
                    this.serverMessage = message;
                    this.serverError = true;
                  }
                },
                ( response) => {
                    console.log(response);
                    this.serverMessage = "服务器开小差儿了";
                    this.serverError = true;
                  }
                );
            },
            checkPassword: function() {
                if (this.entAdminPassword == "") {
                    this.entAdminPasswordRequiredError = true;
                    return false;
                } else {
                    this.entAdminPasswordRequiredError = false;
                }
                return true;
            },
            checkEnterpriseForm:
            function () {
              if (
                this
                  .enterpriseEmailRequiredError ||
                this
                  .enterpriseEmailFormatError ||
                this
                  .enterprisePasswordRequiredError ||
                this
                  .enterprisePasswordFormatError
              ) {
                return false;
              }
              return true;
            },
        },
    }
  );

