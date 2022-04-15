var app = new Vue({
  el: '#app',
  delimiters: ['[[',']]'],
  data: {
      adding: false,
      serverError: false,
      serverMessage: "",
      message: 'Hello Vue!',
      team: '',
      teams: [],
  },
  methods: {
        configTeam: function(item){
            if (item.team != "") {
                window.location.href = "/team/detail/" + item.id;
            }
        },

        configMembers: function(item){
            if (item.team != "") {
                window.location.href = "/teamuser/" + item.id;
            }
        },
        deleteTeam: function(item){
	    if (item.team != "") {
                window.location.href = "/team/delete/" + item.id;
            }
        },
        saveTeam: function(item){
            this.adding = false;
            var data = {
                "name": item.temp_team,
            };

            for (var i = 0; i < this.teams; i++ ){
                var temp = this.teams[i];
                if (item.temp_team == temp.team){
                    return;
                }
            }

            this.$http.post('/api/team/create', data, {emulateJSON: true, credentials : true}).then(response => {
                 var code = response.data.code;
                 if (code == 0) {
                     this.serverMessage = "";
                     this.serverError = false;
                     this.getTeams();
                 } else {
                     var message = response.data.message;
                     this.serverMessage = message;
                     this.serverError = true;
                 }
              }, response => {
                 console.log(response);
                 this.serverMessage = "服务器开小差儿了";
                 this.serverError = true;
            });
        },
        cancleTeam: function(index) {
            this.teams.splice(index, 1);
            this.serverMessage = "";
            this.serverError = false;
        },
        addTeam: function() {
            var team = {
                team:'',
                url: '',
                valid: false,
                type: 0,
              };

              this.teams.push(team)
              this.team = '';
        },
        getTeams: function() {
            var data = {};
            this.$http.get('/api/team/list', data, {emulateJSON: true, credentials : true}).then(response => {
                 var code = response.data.code;
                 if (code == 0) {
                     this.serverMessage = "";
                     this.serverError = false;
                     this.teams = response.data.data.list;
                     console.log(this.teams);
                 } else {
                     var message = response.data.message;
                     this.serverMessage = message;
                     this.serverError = true;
                 }
              }, response => {
                 console.log(response);
                 this.serverMessage = "服务器开小差儿了";
                 this.serverError = true;
            });
        },
   },
    created() {
        this.getTeams();
    }
})


