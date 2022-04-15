package render

import (
	"embed"
	"fmt"
	"html/template"
	"strings"

	"github.com/gin-contrib/multitemplate"
	"github.com/golang/glog"
)

var Render multitemplate.Render
var TemplateBox embed.FS

func InitRender() {
	Render = render(TemplateBox)
}

func BuildTemplate(template_name string, templates []string) (*template.Template, error) {
	temp, err := template.New(template_name).Parse(
		strings.Join(templates, " "))
	return temp, err
}

func getContent(f embed.FS, filename string) string {

	data, err := f.ReadFile(filename)
	if err != nil {
		glog.Error(err)
		return ""
	} else {
		return string(data)
	}
}

func privacyRender(f embed.FS, r *multitemplate.Render) error {

	privacy_templates := []string{
		getContent(f, "templates/base/base.html"),
		getContent(f, "templates/base/search.html"),
		getContent(f, "templates/base/footer.html"),
		getContent(f, "templates/base/track.html"),
		getContent(f, "templates/base/navi.html"),
		getContent(f, "templates/privacy/privacy.html"),
	}

	privacy_tmp, err := BuildTemplate("privacy", privacy_templates)
	if err != nil {
		glog.Error(err)
		return err
	}
	r.Add("Privacy", privacy_tmp)
	return nil
}

func termsRender(f embed.FS, r *multitemplate.Render) error {

	terms_templates := []string{
		getContent(f, "templates/base/base.html"),
		getContent(f, "templates/base/search.html"),
		getContent(f, "templates/base/footer.html"),
		getContent(f, "templates/base/track.html"),
		getContent(f, "templates/base/navi.html"),
		getContent(f, "templates/terms/terms.html"),
	}

	terms_tmp, err := BuildTemplate("terms", terms_templates)
	if err != nil {
		glog.Error(err)
		return err
	}
	r.Add("Terms", terms_tmp)
	return nil
}

func priceRender(f embed.FS, r *multitemplate.Render) error {

	templates := []string{
		getContent(f, "templates/base/base.html"),
		getContent(f, "templates/base/search.html"),
		getContent(f, "templates/base/footer.html"),
		getContent(f, "templates/base/track.html"),
		getContent(f, "templates/base/navi.html"),
		getContent(f, "templates/price/price.html"),
	}

	tmp, err := BuildTemplate("price", templates)
	if err != nil {
		glog.Error(err)
		return err
	}
	r.Add("Price", tmp)
	return nil
}

func billingRender(f embed.FS, r *multitemplate.Render) error {

	templates := []string{
		getContent(f, "templates/base/base.html"),
		getContent(f, "templates/base/search.html"),
		getContent(f, "templates/base/footer.html"),
		getContent(f, "templates/base/track.html"),
		getContent(f, "templates/base/navi.html"),
		getContent(f, "templates/billing/billing.html"),
	}

	tmp, err := BuildTemplate("billing", templates)
	if err != nil {
		glog.Error(err)
		return err
	}
	r.Add("Billing", tmp)
	return nil
}

func teamDeleteRender(f embed.FS, r *multitemplate.Render) error {

	TeamDelete_templates := []string{
		getContent(f, "templates/base/base.html"),
		getContent(f, "templates/base/search.html"),
		getContent(f, "templates/base/footer.html"),
		getContent(f, "templates/base/track.html"),
		getContent(f, "templates/base/navi.html"),
		getContent(f, "templates/TeamDelete/TeamDelete.html"),
	}

	TeamDelete_tmp, err := BuildTemplate("TeamDelete", TeamDelete_templates)
	if err != nil {
		glog.Error(err)
		return err
	}
	r.Add("TeamDelete", TeamDelete_tmp)
	return nil
}

func teamRender(f embed.FS, r *multitemplate.Render) error {

	team_templates := []string{
		getContent(f, "templates/base/base.html"),
		getContent(f, "templates/base/search.html"),
		getContent(f, "templates/base/footer.html"),
		getContent(f, "templates/base/track.html"),
		getContent(f, "templates/base/navi.html"),
		getContent(f, "templates/team/team.html"),
	}

	team_tmp, err := BuildTemplate("team", team_templates)
	if err != nil {
		glog.Error(err)
		return err
	}
	r.Add("Team", team_tmp)
	return nil
}

func memberRender(f embed.FS, r *multitemplate.Render) error {

	member_templates := []string{
		getContent(f, "templates/base/base.html"),
		getContent(f, "templates/base/search.html"),
		getContent(f, "templates/base/footer.html"),
		getContent(f, "templates/base/track.html"),
		getContent(f, "templates/base/navi.html"),
		getContent(f, "templates/member/member.html"),
	}

	member_tmp, err := BuildTemplate("member", member_templates)
	if err != nil {
		glog.Error(err)
		return err
	}
	r.Add("Member", member_tmp)
	return nil
}

func memberDetailRender(f embed.FS, r *multitemplate.Render) error {

	memberDetail_templates := []string{
		getContent(f, "templates/base/base.html"),
		getContent(f, "templates/base/search.html"),
		getContent(f, "templates/base/footer.html"),
		getContent(f, "templates/base/track.html"),
		getContent(f, "templates/base/navi.html"),
		getContent(f, "templates/memberDetail/memberDetail.html"),
	}

	memberDetail_tmp, err := BuildTemplate("memberDetail", memberDetail_templates)
	if err != nil {
		glog.Error(err)
		return err
	}
	r.Add("MemberDetail", memberDetail_tmp)
	return nil
}

func memberDeleteRender(f embed.FS, r *multitemplate.Render) error {

	memberDelete_templates := []string{
		getContent(f, "templates/base/base.html"),
		getContent(f, "templates/base/search.html"),
		getContent(f, "templates/base/footer.html"),
		getContent(f, "templates/base/track.html"),
		getContent(f, "templates/base/navi.html"),
		getContent(f, "templates/memberDelete/memberDelete.html"),
	}

	memberDelete_tmp, err := BuildTemplate("memberDelete", memberDelete_templates)
	if err != nil {
		glog.Error(err)
		return err
	}
	r.Add("MemberDelete", memberDelete_tmp)
	return nil
}

func ticketRender(f embed.FS, r *multitemplate.Render) error {

	ticket_templates := []string{
		getContent(f, "templates/base/base.html"),
		getContent(f, "templates/base/search.html"),
		getContent(f, "templates/base/footer.html"),
		getContent(f, "templates/base/track.html"),
		getContent(f, "templates/base/navi.html"),
		getContent(f, "templates/ticket/ticket.html"),
	}

	ticket_tmp, err := BuildTemplate("ticket", ticket_templates)
	if err != nil {
		glog.Error(err)
		return err
	}
	r.Add("Ticket", ticket_tmp)
	return nil
}

func ticketDetailRender(f embed.FS, r *multitemplate.Render) error {

	ticketDetail_templates := []string{
		getContent(f, "templates/base/base.html"),
		getContent(f, "templates/base/search.html"),
		getContent(f, "templates/base/footer.html"),
		getContent(f, "templates/base/track.html"),
		getContent(f, "templates/base/navi.html"),
		getContent(f, "templates/ticketDetail/ticketDetail.html"),
	}

	ticketDetail_tmp, err := BuildTemplate("ticketDetail", ticketDetail_templates)
	if err != nil {
		glog.Error(err)
		return err
	}
	r.Add("TicketDetail", ticketDetail_tmp)
	return nil
}

func loginRender(f embed.FS, r *multitemplate.Render) error {

	login_templates := []string{
		getContent(f, "templates/base/base.html"),
		getContent(f, "templates/base/search.html"),
		getContent(f, "templates/base/footer.html"),
		getContent(f, "templates/base/track.html"),
		getContent(f, "templates/login/login.html"),
	}

	login_tmp, err := BuildTemplate("login", login_templates)
	if err != nil {
		glog.Error(err)
		return err
	}
	r.Add("Login", login_tmp)
	return nil
}

func signupRender(f embed.FS, r *multitemplate.Render) error {
	signup_templates := []string{
		getContent(f, "templates/base/base.html"),
		getContent(f, "templates/base/search.html"),
		getContent(f, "templates/base/footer.html"),
		getContent(f, "templates/base/track.html"),
		getContent(f, "templates/signup/signup.html"),
	}

	signup_tmp, err := BuildTemplate("signup", signup_templates)
	if err != nil {
		glog.Error(err)
		return err
	}
	r.Add("Signup", signup_tmp)
	return nil
}

func forgotPasswordRender(f embed.FS, r *multitemplate.Render) error {
	forgotPassword_templates := []string{
		getContent(f, "templates/base/base.html"),
		getContent(f, "templates/base/search.html"),
		getContent(f, "templates/base/footer.html"),
		getContent(f, "templates/base/track.html"),
		getContent(f, "templates/forgotPassword/forgotPassword.html"),
	}

	forgotpassword_tmp, err := BuildTemplate("forgotpassword", forgotPassword_templates)
	if err != nil {
		fmt.Println(err)
		glog.Error(err)
		return err
	}
	r.Add("ForgotPassword", forgotpassword_tmp)
	return nil
}

func homeRender(f embed.FS, r *multitemplate.Render) error {
	home_templates := []string{
		getContent(f, "templates/base/base.html"),
		getContent(f, "templates/base/search.html"),
		getContent(f, "templates/base/navi.html"),
		getContent(f, "templates/base/footer.html"),
		getContent(f, "templates/base/track.html"),
		getContent(f, "templates/home/home.html"),
	}

	home_tmp, err := BuildTemplate("Home", home_templates)
	if err != nil {
		glog.Error(err)
		return err
	}

	r.Add("Home", home_tmp)
	return nil
}

func aboutRender(f embed.FS, r *multitemplate.Render) error {
	about_home_templates := []string{
		getContent(f, "templates/base/base.html"),
		getContent(f, "templates/base/search.html"),
		getContent(f, "templates/base/navi.html"),
		getContent(f, "templates/base/footer.html"),
		getContent(f, "templates/base/track.html"),
		getContent(f, "templates/about/about.html"),
	}
	about_home_tmp, err := BuildTemplate("about", about_home_templates)
	if err != nil {
		glog.Error(err)
		return err
	}
	r.Add("About", about_home_tmp)
	return nil
}

func docApiRender(f embed.FS, r *multitemplate.Render) error {
	doc_api_home_templates := []string{
		getContent(f, "templates/base/base.html"),
		getContent(f, "templates/base/search.html"),
		getContent(f, "templates/base/navi.html"),
		getContent(f, "templates/base/footer.html"),
		getContent(f, "templates/base/track.html"),
		getContent(f, "templates/doc/api.html"),
	}
	doc_api_home_tmp, err := BuildTemplate("APIDoc", doc_api_home_templates)
	if err != nil {
		glog.Error(err)
		return err
	}
	r.Add("APIDoc", doc_api_home_tmp)
	return nil
}

func errorRender(f embed.FS, r *multitemplate.Render) error {
	error_templates := []string{
		getContent(f, "templates/base/base.html"),
		getContent(f, "templates/base/search.html"),
		getContent(f, "templates/base/navi.html"),
		getContent(f, "templates/base/footer.html"),
		getContent(f, "templates/base/track.html"),
		getContent(f, "templates/error/error.html"),
	}
	error_tmp, err := BuildTemplate("error", error_templates)
	if err != nil {
		glog.Error(err)
		return err
	}
	r.Add("Error", error_tmp)
	return nil
}

func signupMailRender(f embed.FS, r *multitemplate.Render) error {
	error_templates := []string{
		getContent(f, "templates/email/signup.html"),
	}
	error_tmp, err := BuildTemplate("signupEmail", error_templates)
	if err != nil {
		glog.Error(err)
		return err
	}
	r.Add("SignupEmail", error_tmp)
	return nil
}

func forgotPasswordMailRender(f embed.FS, r *multitemplate.Render) error {
	error_templates := []string{
		getContent(f, "templates/email/forgotpassword.html"),
	}
	error_tmp, err := BuildTemplate("forgotpasswordEmail", error_templates)
	if err != nil {
		glog.Error(err)
		return err
	}
	r.Add("ForgotPasswordEmail", error_tmp)
	return nil
}

func ticketReplyMailRender(f embed.FS, r *multitemplate.Render) error {
	error_templates := []string{
		getContent(f, "templates/email/ticketreply.html"),
	}
	error_tmp, err := BuildTemplate("ticketReply", error_templates)
	if err != nil {
		glog.Error(err)
		return err
	}
	r.Add("TicketReplyEmail", error_tmp)
	return nil
}

func render(f embed.FS) multitemplate.Render {

	r := multitemplate.New()

	homeRender(f, &r)
	memberRender(f, &r)
	memberDeleteRender(f, &r)
	memberDetailRender(f, &r)
	teamRender(f, &r)
	teamDeleteRender(f, &r)
	aboutRender(f, &r)
	termsRender(f, &r)
	privacyRender(f, &r)
	errorRender(f, &r)

	loginRender(f, &r)
	signupRender(f, &r)
	forgotPasswordRender(f, &r)

	signupMailRender(f, &r)
	forgotPasswordMailRender(f, &r)
	return r
}
