package parse

import (
	"testing"
)

func TestFindPrice(t *testing.T) {
	failed := false
	testPairs := map[string]int{
		"<i>$<a>45</a></i>": 45,
		"<p class='BpkText_bpk-text__2NHsO BpkText_bpk-text--sm__345aT FqsTabs_fqsTabWithSparkleText__mE9wc FqsTabs_fqsTabWithSparkleTextSelected__1T7_o'>Best</p><div class='FqsTabs_fqsTabWithSparklePrice__1-UyJ'><div class='Price_mainPriceContainer__1dqsw'><span class='BpkText_bpk-text__2NHsO BpkText_bpk-text--lg__3vAKN BpkText_bpk-text--bold__4yauk'>$1,465</span></div></div><p class='BpkText_bpk-text__2NHsO BpkText_bpk-text--sm__345aT FqsTabs_fqsTabWithSparkleText__mE9wc FqsTabs_fqsTabWithSparkleDuration__1ugSs FqsTabs_fqsTabWithSparkleTextSelected__1T7_o'>20h 52m (average)</p>":                                                                                                                                                                                                                                                                                                                                                                        1465,
		"<span class='BpkText_bpk-text__2NHsO BpkText_bpk-text--sm__345aT DealsFrom_dealsText__C6CoA'>1 deal</span><div class='><div class='Price_mainPriceContainer__1dqsw'><span class='BpkText_bpk-text__2NHsO BpkText_bpk-text--lg__3vAKN BpkText_bpk-text--bold__4yauk'>$1,465</span></div></div><button type='button' class='BpkButton_bpk-button__32HxR TicketStub_ctaButton__2ctIF'>Select&nbsp;<span style='line-height: 1.125rem; display: inline-block; margin-top: 0.1875rem; vertical-align: top;'><svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' width='18' height='18' style='width: 1.125rem; height: 1.125rem;' class='BpkIcon_bpk-icon--rtl-support__2-_zb' fill='white'><path d='M14.4 19.5l5.7-5.3c.4-.4.7-.9.8-1.5.1-.3.1-.5.1-.7s0-.4-.1-.6c-.1-.6-.4-1.1-.8-1.5l-5.7-5.3c-.8-.8-2.1-.7-2.8.1-.8.8-.7 2.1.1 2.8l2.7 2.5H5c-1.1 0-2 .9-2 2s.9 2 2 2h9.4l-2.7 2.5c-.5.4-.7 1-.7 1.5s.2 1 .5 1.4c.8.8 2.1.8 2.9.1z'></path></svg></span></button>": 1465,
	}
	for input, expectation := range testPairs {
		got := FindPrice(input)

		if got != expectation {
			t.Error(got, expectation)
			failed = true
		}
	}
	if failed {
		t.Fail()
	}
}
