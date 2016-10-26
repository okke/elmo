package elmo

import "testing"

func TestLoaderLoadsFromWorkingDir(t *testing.T) {

	loader := NewLoader(NewGlobalContext(), []string{"../examples/basics"})

	notFound := loader.Load("chipotles")
	if notFound.Type() != TypeError {
		t.Error("why did we find chipotles?")
	}

	value := loader.Load("simple-functions")

	if value.Type() != TypeDictionary {
		t.Errorf("loading should result in a dictionary, found %v", value)
	}
}

func TestLoaderDoesNotLoadFileWithError(t *testing.T) {

	loader := NewLoader(NewGlobalContext(), []string{"./loader_testdata"})

	value := loader.Load("undefined-function")

	if value.Type() != TypeError {
		t.Errorf("could load undefined-functions: %v", value)
	}

	if value.String() != "error at ./loader_testdata/undefined-function.mo at line 2: call to undefined \"what\"" {
		t.Errorf("expected a different value, found %v", value)
	}

}

func TestLoaderLoadFileThatLoadsAFile(t *testing.T) {

	loader := NewLoader(NewGlobalContext(), []string{"./loader_testdata"})

	value := loader.Load("use_load")

	if value == nil || value.Type() != TypeDictionary {
		t.Errorf("could not load use_load: %v", value)
	}

}

func TestLoaderDoesNotLoadFileWhichCanNotLoadOtherFile(t *testing.T) {

	loader := NewLoader(NewGlobalContext(), []string{"./loader_testdata"})

	value := loader.Load("undefined-script")

	if value.Type() != TypeError {
		t.Errorf("could load undefined-script: %v", value)
	}

	if value.String() != "error at ./loader_testdata/undefined-script.mo at line 1: mixin can only mix in dictionaries, not error: could not find szechuan" {
		t.Errorf("expected a different value, found %v", value)
	}

}
