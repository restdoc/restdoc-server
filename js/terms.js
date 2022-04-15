var app = new Vue({
  el: '#app',
  delimiters: ['[[',']]'],
  data: {
      serverError: false,
      serverMessage: "",
      message: 'Hello Vue!',
      email: '',
      emailRequiredError: false,
      emailFormatError: false,
      emailPattern : /^([A-Za-z0-9_\-\.\u4e00-\u9fa5])+\@([A-Za-z0-9_\-\.])+\.([A-Za-z]{2,8})$/,
      password: '',
      passwordRequiredError: false,
  },
    methods: {
    }
})

