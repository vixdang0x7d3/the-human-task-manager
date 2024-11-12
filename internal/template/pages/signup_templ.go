// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.793
package pages

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import "github.com/vixdang0x7d3/the-human-task-manager/internal/template"

func Signup() templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Var2 := templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
			templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
			templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
			if !templ_7745c5c3_IsBuffer {
				defer func() {
					templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
					if templ_7745c5c3_Err == nil {
						templ_7745c5c3_Err = templ_7745c5c3_BufErr
					}
				}()
			}
			ctx = templ.InitializeContext(ctx)
			templ_7745c5c3_Err = template.Head().Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" <body><div class=\"flex items-center justify-center py-10\"><div class=\"card bg-base-200 w-120 shadow-xl\"><div class=\"card-body\" class=\"card-body items-center text-center\"><h2 class=\"card-title pb-4\">Sign up new account</h2><form hx-post=\"/v1/users\" hx-target=\"#card-body\"><div class=\"flex flex-col gap-2\"><label class=\"input input-bordered flex items-center gap-2\">Username: <input name=\"username\" id=\"username\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if true {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" required")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				return templ_7745c5c3_Err
			})
			templ_7745c5c3_Err = template.Head().Render(templ.WithChildren(ctx, templ_7745c5c3_Var3), templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" <body>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = template.Navbar().Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"flex items-center justify-center py-10\"><div class=\"card bg-base-200 w-120 shadow-xl\" hx-ext=\"response-targets\"><div id=\"card-body\" class=\"card-body items-center text-center\"><h2 class=\"card-title pb-4\">Sign up new account</h2><form hx-post=\"/v1/users\" hx-target=\"#card-body\" hx-target-400=\"#message\"><div class=\"flex flex-col gap-2\"><label class=\"input input-bordered flex items-center gap-2\">Username: <input type=\"text\" name=\"username\" id=\"username\" class=\"grow\" placeholder=\"taskman\"></label> <label class=\"input input-bordered flex items-center gap-2\">Email: <input type=\"email\" name=\"email\" id=\"email\" class=\"grow\" placeholder=\"htm@site.com\"></label> <label class=\"input input-bordered flex items-center gap-2\">First Name: <input type=\"text\" name=\"first_name\" id=\"first_name\" class=\"grow\" placeholder=\"Human\"></label> <label class=\"input input-bordered flex items-center gap-2\">Last Name: <input type=\"text\" name=\"last_name\" id=\"last_name\" class=\"grow\" placeholder=\"Task Manager\"></label> <label class=\"input input-bordered flex items-center gap-2\">Password: <input type=\"password\" name=\"password\" id=\"password\" class=\"grow\" placeholder=\"********\"></label></div><div id=\"message\"></div><div class=\"card-actions justify-end pt-4\"><input type=\"submit\" class=\"btn btn-primary\" value=\"Submit\"> <a href=\"#!\" class=\"btn btn-secondary\">Cancel</a></div></form></div></div></div></body>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			return templ_7745c5c3_Err
		})
		templ_7745c5c3_Err = template.Boilerplate().Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate
