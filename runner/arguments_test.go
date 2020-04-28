package runner

import "testing"

func TestNextArgOnEmptyArray(t *testing.T) {

	var args cliArgs = []string{}

	first, more := args.next()

	if first != "" {
		t.Error("did not expect a first argument")
	}

	if len(more) != 0 {
		t.Error("did not expect a more argument")
	}

}

func TestNextArgOnArrayWithOneElem(t *testing.T) {

	var args cliArgs = []string{"chipotle"}

	first, more := args.next()

	if first != "chipotle" {
		t.Error("expected chipotle as first argument")
	}

	if len(more) != 0 {
		t.Error("did not expect a more argument")
	}

}

func TestNextArgOnArrayWithTwoElems(t *testing.T) {

	var args cliArgs = []string{"chipotle", "jalapeno"}

	first, more := args.next()

	if first != "chipotle" {
		t.Error("expected chipotle as first argument")
	}

	if len(more) != 1 {
		t.Error("expected one more argument")
	}

}

func TestPutBack(t *testing.T) {

	var args cliArgs = []string{"chipotle", "jalapeno"}

	args = args.putBack("habanero")

	first, more := args.next()

	if first != "habanero" {
		t.Error("expected chipotle as first argument")
	}

	if len(more) != 2 {
		t.Error("expected otwo more arguments")
	}

}

func TestParseFlagAsOneArg(t *testing.T) {

	mapping := make(map[string]string, 0)
	var args cliArgs = []string{"-pepper=jalapeno"}

	parseFlags(args, func(name, value string) {
		mapping[name] = value
	})

	value, found := mapping["pepper"]
	if !found {
		t.Errorf("expected a pepper flag in %v", mapping)
	} else if value != "jalapeno" {
		t.Errorf("expected a jalapeno pepper, not %s in %v", value, mapping)
	}

}

func TestParseFlagAsSeparateArg(t *testing.T) {

	mapping := make(map[string]string, 0)
	var args cliArgs = []string{"-pepper", "-sauce"}

	parseFlags(args, func(name, value string) {
		mapping[name] = value
	})

	if value, found := mapping["pepper"]; !found {
		t.Errorf("expected a pepper flag in %v", mapping)
	} else if value != "true" {
		t.Errorf("expected a true pepper, not %s", value)
	}

	if value, found := mapping["sauce"]; !found {
		t.Errorf("expected a sauce flag in %v", mapping)
	} else if value != "true" {
		t.Errorf("expected a true sauce, not %s", value)
	}

}

func TestParseArguments(t *testing.T) {

	testSetter := newRunnerArgs()
	parseArguments([]string{"-pepper=jalapeno"}, testSetter)
	if _, found := testSetter.elmoFlags["pepper"]; !found {
		t.Error("no pepper arg found")
	}

	testSetter = newRunnerArgs()
	parseArguments([]string{"-pepper=jalapeno", "chipotle.mo"}, testSetter)
	if testSetter.elmoFile != "chipotle.mo" {
		t.Error("no elmo file found")
	}

	testSetter = newRunnerArgs()
	parseArguments([]string{"-pepper=jalapeno", "chipotle.mo", "-pepper=habanero"}, testSetter)
	if _, found := testSetter.elmoFlags["pepper"]; !found {
		t.Error("no pepper arg found")
	}
	if testSetter.elmoFile != "chipotle.mo" {
		t.Error("no elmo file found")
	}
	if _, found := testSetter.userFlags["pepper"]; !found {
		t.Error("no pepper arg in user flags found")
	}

	testSetter = newRunnerArgs()
	parseArguments([]string{"-repl", "-debug", "chipotle.mo", "-pepper"}, testSetter)
	if _, found := testSetter.elmoFlags["repl"]; !found {
		t.Error("no repl arg found")
	}
	if _, found := testSetter.elmoFlags["debug"]; !found {
		t.Error("no debug arg found")
	}
	if testSetter.elmoFile != "chipotle.mo" {
		t.Error("no elmo file found")
	}
	if _, found := testSetter.userFlags["pepper"]; !found {
		t.Error("no debug arg found")
	}
}
