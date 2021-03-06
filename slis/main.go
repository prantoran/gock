package main

import (
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
	"magic.pathao.com/pinku/socktest/conf"
)

var userTk = []string{
	"kPiBSyNStlRsGB7OigtMKVXjc7iTV6M4VLswXzoY",
	"8NcgOlpHnSEd8dnKjW9HEOVqFDr3ObykB4kBSZH2",
	"n1UsEHLGhRjl8rgnpv7FHqPjASWpdyIw1zai4xPs",
	"9WN0iCzfUW9rGwllbpepDBZV7I6DPSnQ4geD0fLZ",
	"ZTrcqVWfk8dynWoxW5YyFHNju5pMJspMSGIstqhk",
	"NC4CvxhVhQX3UliTeE8XP7AAZUaEGAuT9vhhnYD0",
	"teHAcGrIImUDIYZv9g5duW0Q5uxajgounBiABWtM",
	"mqzroLXOj3potYorMzG29xVy0wA2Unrstgg9Vi79",
	"luBPtccuz0rbNQYf3iBQjuVWXJxSqYIs0doYwWnO",
	"sg9yhKelBXW2daB0YODbjuxTLbICi7Xn7EGOx7gJ",
	"zRuzZoBdJavugLRxdOAhMtIE7iSZdZbIba2MoLE1",
	"DY4U3vBgEknflgwUzRuynUUazdSzpbiXNhrP272Q",
	"vtBMF9tfwDmERWtCkpHLaZtFJYL8qxFtCj76JvQ8",
	"uVns9fsbVq0TBS00L3Tmfur381qWcuF7lsAftAvn",
	"BbjEj9Sm6ySibkQRQGHfBxCxjI0dnb6SilH330E7",
	"owgCEv7XsOoTDbztviV3qw9n7RHh1gio98srltPP",
	"hxqmr6N4xES8Gtuin17ckX2Nu0tYb9koW8zMs56x",
	"aQbjUxUayaEPRPSOwbjMs4NkvX2qxqECBi9P2PH9",
	"giP1XKh1qjj8htUWHJL6DJ4kp4LPoCn5gELEkYsE",
	"QSSNh1AmMAwvmrdaOh4hmIsKKlGh1DXPxr4t71Rt",
	"PnR5KlqzK3XnychmI3tsk2lbRh8639ARSMd3Groi",
	"eoZQfUSvwphkuEhAS8P6B68u56oSYDoUkcZ1KV1J",
	"MNUrc635PNqUxhr6B4bOEo46dSl7oGidmdMOK80B",
	"dkGhxakgpXzIjKDrQJsTRye3H8nU4IwjhXXvEymn",
	"wvC8wrxJMmfl9MDI3EmIjRFxzUunUxHyu0o366QH",
	"jp8fxymafZJnlNUtkih9EJ2R7hQPEWae84MKIELo",
	"tN3qENVHWncQT0Mn05MvZalYysrHhGu3KOrA6B8b",
	"jEnsgsrp9CaU7ud90MZZW6FLOeFbJgV1G6RJ9Pcb",
	"YXnLSEd8UtZooX34LXVlfdwsCPS667Ub8KhTAil9",
	"LxOTQsB0t8E6CiAQoKZc2vnrHHv4PPAFhnJsW9n9",
	"fNDiv9XjVE93EVeeEtQQXXkxdLWQ0XZanROcDYeW",
	"nbTf71A9aY1jgIyMuks8RA2T5L12RkRrpnNnKqcz",
	"hTisPoPaIpjCqOXVteTvMvofGB9SwY7AoZVWhjwr",
	"5WD5lMegwHTHb26wOu8bgP0QYCh9vxjTxo6h82W7",
	"EVncIKD7vIyXVoeIiWWmFr5J6pMylXzN5Dc0Y1Qh",
	"UYR464VxxMH0lNETKDyTs7MS4bfghqGIEzGDo9nU",
	"JhCHhy8S8wyaqMnKth9e6sDpq7uKM0cPMKmC4TXS",
	"GDmoZqx98pyV9hIvl74pGpdtEkB6HXZ15NzaCvsu",
	"cKTQRhyLXQt1OR4UccCDQpzRc3No4SwUbV7skCnP",
	"clKAL3E80NpmtXKMtYTY7vgYcJknmUZ4TJAfUqQS",
	"iN2zdyLRkkXwvg2ZNrIWBbJlwS8IjZyfgTc466YQ",
	"3s4vKwRLpJ9RWEu0LmB6aGKHe7DSSRDLiCQPD3Eu",
	"TbqoBow3YyrnvFLAR12g0FgdD3ZmH2q77gIOe2fM",
	"Ujfc8Wf2NPFj9Tfx3WUTZshhsr42H54kaHk4fxl6",
	"M2WX6ktuYgXU6XcbIZeWI5lwIVyBagU4cZGRM5bM",
	"3W5WVZZzSBLge3XzmILwuiYi6LvehcrdqELitavA",
	"VSd0q5LKUBpt3OjfCi2rrnheyKxA4cOaqf5NDzpK",
	"cQsfimREp6K4fXZl3RM4vERkBHA3kYZeYHqSfEDv",
	"2afoEu5N7SkLwfAnPo1BP6FexzfO1Y2BqH9N5hiB",
	"4sF886plTGVshJIIhidLsYI1LS9VQXRuCI2YxA6A",
	"p7wMX86t0QKpRPcDUvqzFeg9WWFdDHOstk24YkNx",
	"Mha3Kg5K57n5zOrdxP1d2vDMCGNxfqNg4iA9xk4d",
	"Dnqw6XR8FkAEEn1SyWScEIhkS4IUjaqhtXdLlFdG",
	"GfVyyWpRw5eyuYEev9UkmVHj1bVFWA8H7z9AkIzZ",
	"gJHkLLpdnfeQPaV9RdxYYsDHBH1kmTOEsNtasVoW",
	"y2HMrom8OD4VdZJdP0llsoYyKpF0KyXjz871ghB9",
	"CwLn1ARfKtxsj2j457gYlbSRISkddLVuv9BLKhdF",
	"NQ0Hh234S0Wvsi10kMRuBY7ZPXCzo8JSc1bF811q",
	"zIwVDpIrvWV5AjXe1ftsJ2E6vfjc7qT6ZwweERiz",
	"23YyHmMMJNhx4isLGAjsP5aMSxDGmVS4GVUfyTvc",
	"Gae6RBNKGOltz9B2CXQdUv3odaOlqGTljYFTxm2B",
	"DRbMbMzPmCf9p7gBLtgALnt0Dn41Ck1p2CYrO8jm",
	"Veqd8ux9nRoBODnhDFtcqfa3XihYo52iFNtgXGy5",
	"4wV3GjRkutc4J0i97qEvM3ESe2kK6ZPMDgAxMHYL",
	"p0w5nqxqt0dLEu2EbdPZ0JHLJAiL21F9aa0FLF7s",
	"e3xP2uBr0RTesiLBbuxxqnto4b5UbFsg54vOgIRC",
	"tmELf94GJI62yngLOFMrGQmRnI90udA6g1CKltuy",
	"BpfkRWcl014dpjdkhXRg0DQj0bKwgmVuj0SjeWXL",
	"1omWEkK9gfmbaKZ5zmcAQG8LmwqYfGGcnebCqK1M",
	"ZRG5dnOZ3P4ucgW5Gut3ufsDDdOkka8v6uDXMFnD",
	"YYxai2Ve6M4Fs2Xt53Fdsz2hvcMJtVjhxWyrLFyC",
	"G0g4kszUk3gigTyotzT8eGH9Nl2LG0uAXsDH7Gh5",
	"i164Ra9x6ug4HosSflQ3GTjnTrB2Cb9qSAeQ7epa",
	"6vnpju56lNlnPzNK6uUtjYxcLzkha1JdM2WDv9pK",
	"IL37ddrUxX2D1Wq4yBT3HIRKHBv3O1ptLvzVvqoW",
	"efUXGRDAFcuvtInqDJ1pyNbT9JxLSHWy2QhKJF6p",
	"vLqoqsQpNcHavfSHXfssvzVWvCeVt7UPzzHlYxKQ",
	"UetZdXG0OszKoMvFHtJ4pymVOqAklTOJdta1GHIw",
	"ja2bAgnEKAQkOF1PtAfJGPC7DkCUbIA8Cp4A2Uw6",
	"Nt7X4g7F6Dntlg6Ym0U5nvM2rLOMnOddZHe8bBj4",
	"lT6Azg8c2NKpGLT0iP6CeDw8oZOIqUrfqpoXRTI6",
	"HOsI5Q6Gv3DKbRTXgityyglzADlMpNXdIQd6jced",
	"NPgJjRonHfhyX6yTILFUYr8fttOlwLzSpzn14axR",
	"TplhhNfK5QfiB4Yo92EsMl4mVml1WPOxL3m1IFpj",
	"GNJlfCBcVXCGTKIIMrI9hVG8CEVu6ujwlLRfmBpb",
	"rp2HzE5ItyMmlsRHsPvApltxLSd8BrBNL0u2vKWq",
	"ASHQWGxomks4EdI4DUffeVn4enWmRPoayUx7rPAp",
	"ouw5q1vqHMSPxnr9qtEyHoNrkdhBF9CnuEIjMW3i",
	"qfVPkH1BoIBeEkL0pr3N6VoMp5gthlkmjHUScZoH",
	"ENOC6swRoOv7v1g74qmOHetBWoPFlodt6LRWWjPT",
	"VbM9En5MiwG23HlyeA7CZhLC0WXiyfTxELO6lZSu",
	"XU86nc6pNhIIwSvVHiFOX0I5RLpwb92ft9u3ivLm",
	"ItRqiwRyaeqHewC5Te0gLnjMgcixdHQrK2ji2A2T",
	"aKgQT0n17AgOLP78DgZMdXs0z9I59bx3Sq5z6BU7",
	"GbBK9VYHwkpuPl6SMXMKQ8pibi100s3hKZbwfKBH",
	"0Sx0rxItDxtr7JhMElgLEjJIUSeXaNPUsX2ZSrrb",
	"G2hq1sjjoAhqt45PVQ6cuOWq4O2plgKyXBOtOuxm",
	"fTh3QnulAdFzhBjq47luPe0xbsUwdBrIW8mDthyj",
	"u1g2KQCxcngpVIWm9OVtiynirL6QcorcIbfC0NC8",
	"caHjsMUmzs7t9RdmZVB1jXGYXNlVvXfU90ksH62V",
}

func listener(host string, token string) error {
	u := url.URL{Scheme: "ws", Host: host, Path: ""}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}

	c.WriteJSON(map[string]interface{}{
		"token": token,
		"type":  "register",
	})

	go func(conn *websocket.Conn) {
		defer func() {
			conn.Close()
		}()

		for {
			msg := make(map[string]interface{})
			err := websocket.ReadJSON(conn, &msg)
			if err != nil {
				log.Println("listener:", err)
				return
			}

			log.Println("get message", msg)
		}
	}(c)

	return nil
}

func main() {
	if err := conf.Load(); err != nil {
		log.Println("Could not load config:", err)
		return
	}
	fmt.Println("conf.Addr:", conf.Addr)
	for i := range userTk {
		log.Println("i:", i, " ", listener(conf.Addr, userTk[i]))
		// log.Println("i:", i, " ", listener("35.201.176.61:8081", userTk[i]))
		break
	}

	forever := make(chan bool)
	<-forever
}
