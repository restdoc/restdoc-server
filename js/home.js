var app = new Vue({
  el: '#app',
  delimiters: ['[[',']]'],
  data: {
      adding: false,
      serverError: false,
      serverMessage: "",
      message: 'Hello Vue!',
      domain: '',
      selected: 0,
      content: '',
      list: [
      ],
  },
    computed: {
        classObject: function() {
            return {
                is_active: this.current_tab == this.tab,
            }
        }
    },
  methods: {
        select: function(index) {
          this.selected = index;
          this.content = this.list[index]['content'];
          Prism.highlightAll();
        },
        start: function(){
            window.location.href = "/signup";
        },
        configDomain: function(item){
        },
        deleteDomain: function(item){
        },
        saveDomain: function(item){
        },
        cancleDomain: function(index) {
        },
        addDomain: function() {
        },
        getDomains: function() {
        },
    },
    created() {
    }
})
