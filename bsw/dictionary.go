package bsw

import (
	"errors"
	"github.com/miekg/dns"
)

func GetWildCard(domain string, serverAddr string) string {
	var fqdn = "youmustconstructmorepylons." + domain
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(fqdn), dns.TypeA)
	in, err := dns.Exchange(m, serverAddr+":53")
	if err != nil {
		return ""
	}
	if len(in.Answer) < 1 {
		return ""
	}
	if a, ok := in.Answer[0].(*dns.A); ok {
		return a.A.String()
	} else {
		return ""
	}
}

func GetWildCard6(domain string, serverAddr string) string {
	var fqdn = "youmustconstructmorepylons." + domain
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(fqdn), dns.TypeAAAA)
	in, err := dns.Exchange(m, serverAddr+":53")
	if err != nil {
		return ""
	}
	if len(in.Answer) < 1 {
		return ""
	}
	if a, ok := in.Answer[0].(*dns.AAAA); ok {
		return a.AAAA.String()
	} else {
		return ""
	}
}

func Dictionary(domain string, subname string, blacklist string, serverAddr string) ([]Result, error) {
	results := make([]Result, 0)
	var fqdn = subname + "." + domain
	ip, err := LookupName(fqdn, serverAddr)
	if err != nil {
		cfqdn, err := LookupCname(fqdn, serverAddr)
		if err != nil {
			return results, err
		}
		ip, err = LookupName(cfqdn, serverAddr)
		if err != nil {
			return results, err
		}
		if ip == blacklist {
			return results, errors.New("Returned IP in blackslist")
		}
		results = append(results, Result{Source: "Dictionary-CNAME", IP: ip, Hostname: fqdn}, Result{Source: "Dictionary-CNAME", IP: ip, Hostname: cfqdn})
		return results, nil
	}
	if ip == blacklist {
		return results, errors.New("Returned IP in blacklist")
	}
	results = append(results, Result{Source: "Dictionary", IP: ip, Hostname: fqdn})
	return results, nil
}

func Dictionary6(domain string, subname string, blacklist string, serverAddr string) ([]Result, error) {
	results := make([]Result, 1)
	var fqdn = subname + "." + domain
	ip, err := LookupName6(fqdn, serverAddr)
	if err != nil {
		return results, err
	}
	if ip == blacklist {
		return results, errors.New("Returned IP in blacklist")
	}
	results[0] = Result{Source: "Dictionary IPv6", IP: ip, Hostname: fqdn}
	return results, nil
}
