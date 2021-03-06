package api

import (
	gohttp "net/http"

	"github.com/gorilla/pat"

	"github.com/doctolib/MailHog/pkg/config"
)

func CreateAPI(conf *config.Config, r gohttp.Handler) {
	apiv1 := createAPIv1(conf, r.(*pat.Router))
	apiv2 := createAPIv2(conf, r.(*pat.Router))

	go func() {
		for msg := range conf.MessageChan {
			apiv1.messageChan <- msg
			apiv2.messageChan <- msg
		}
	}()
}
