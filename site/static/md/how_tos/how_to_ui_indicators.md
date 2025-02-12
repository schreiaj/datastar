# UI Indicators

## Intro 

UI Indicators are a way to show the user that something is happening. A simple example of this is a spinner to represent an async action is in progress. 

## Goal

We're going to explore a few different approaches to UI Indicators. 

## Simple Disable Button

Perhaps, the simplest way to indicate the UI is awaiting an async response is to simply disable the button. You see this frequently in forms. 

Datastar simplifies this with the `data-indicator` attribute. It accepts the name of a signal that will be set to `true` when a fetch request is made by the same element.

By using the `data-attr` attribute we can set the disabled attribute of the button to the value of the signal, in this case `$saving`.
```
<div data-signals="{saving: false}">
    <button
    class="btn btn-primary"
     data-attr-disabled="$saving"
     data-on-click="@post('/how_tos/ui-indicators/save')"
     data-indicator-saving>
        Save
    </button>
</div>
```

### Backend
There is nothing special about this, it is a  simple HTTP handler. For demo purposes we've added a 2s delay to the response.

### Demo

<div data-signals="{saving: false}">
    <button
    class="btn btn-primary"
     data-attr-disabled="$saving"
     data-on-click="@post('/how_tos/ui-indicators/save')"
     data-indicator-saving>
        Save
    </button>
</div>

Clicking on  this button will send a POST request to the backend. While this request is in flight the value of the `$saving` signal will be set to `true` and the button will be disabled.

## Spinner

We can use a very similar approach to load a spinner. 

```html
<style>
    .loader {
        width: 1em;
        height: 1em;
        border: 2px solid #FFF;
        border-bottom-color: transparent;
        border-radius: 50%;
        display: inline-block;
        box-sizing: border-box;
        animation: rotation 2s linear infinite;
    }

    @keyframes rotation {
        0% {
            transform: rotate(0deg);
        }
        100% {
            transform: rotate(360deg);
        }
    }
</style>

<div data-signals="{saving: false}">
    <button
            class="btn btn-primary"
            data-attr-disabled="$saving"
            data-on-click="@post('/how_tos/ui-indicators/save')"
            data-indicator-saving>
        Save
    </button>
    <span data-class="{loader: saving}"></span>
</div>
```

### Backend
Similar to the previous example, there is nothing special about this, it is a  simple HTTP handler (in fact it is the same handler).


### Demo

<style>
    .loader {
        width: 1em;
        height: 1em;
        border: 2px solid #FFF;
        border-bottom-color: transparent;
        border-radius: 50%;
        display: inline-block;
        box-sizing: border-box;
        animation: rotation 2s linear infinite;
    }

    @keyframes rotation {
        0% {
            transform: rotate(0deg);
        }
        100% {
            transform: rotate(360deg);
        }
    }
</style>

<div data-signals="{saving: false}">
    <button
        class="btn btn-primary"
        data-attr-disabled="$saving"
        data-on-click="@post('/how_tos/ui-indicators/save')"
        data-indicator-saving>
            Save 
    </button>
    <span data-class="{loader: $saving}"></span>
</div>

Much like the previous example, we are using the `data-indicator` attribute to update the `$saving` signal. This time we've added in a `data-class` attribute to a span that will show a CSS spinner in addition to disabling the Save button. 

## Optimistic Update

Datastar is so fast at updating the dom that you shouldn't need to do client side optimistic updates (even with high network latency). In the most extreme latency situations adding a a simple indicator can be enough to keep the UI feeling snappy.

The other concern with optimistic updates is that they can mislead a user into thinking an action has been completed when it hasn't. Imagine thinking you have successfully booked a plane ticket, only to find out later it was an optimistic update and your booking did not go through.

If you still want to implement the optimistic update pattern in Datastar despite these shortcomings. Here's how you can do it.

```html
<div>
    <div data-signals="{saving: false, name:''}">
        <div  class="text-primary">
            <input  id="name" data-attr-disabled="$saving" type="text" data-bind="name">
            <button class="flex-1 btn btn-primary" data-indicator-fetching data-on-click="@post('/how_tos/ui-indicators/optimistic')" data-attr-disabled="$saving" >
                Save
            </button>
        </div>
        <div id="list">
            <div id="optimistic" class="text-primary" data-show="$saving" data-text="$name"></div>
        </div>
    </div>
</div>
```

### Backend
This backend is a bit more complicated than the other ones. It requires updating some signals and using different merge modes to place DOM content in specific locations.

```go
func(w http.ResponseWriter, r *http.Request) {
        signals := &OptimisticUpdateSignals{}
        datastar.ReadSignals(r, signals)
        sse := datastar.NewSSE(w, r)
        time.Sleep(2 * time.Second)
        sse.MergeSignals([]byte("{name:''}"))
        sse.MergeFragmentTempl(howTosUiIndicatorsUpdateSucess(signals.Name), datastar.WithSelector("#optimistic"), datastar.WithMergeAfter())
    }
```

The important part is to update the `$name` signal and merge the fragment after the #optimistic element. 

When the request is in flight the `$fetching` signal will be set to true. This will disable the button and input using the same approach as other examples here. Because `$fetching` is true the `data-show` attribute will be set to true and the `#optimistic` div will be shown with the content of the `$name` signal which was submitted to the backend. When the backend returns it resets `$name` and `$fetching` is set to false hiding the `#optimistic` div. Then the backend sends a fragment to be inserted as the next sibling element of the `#optimistic` div.

### Demo

<div>
  <div data-signals="{fetching: false, name:''}">
    <div  class="text-primary">
        <input id="name" data-attr-disabled="$fetching" type="text" data-bind="name">
        <button class="flex-1 btn btn-primary" data-indicator-fetching data-on-click="@post('/how_tos/ui-indicators/optimistic')" data-attr-disabled="$fetching" >
            Save
        </button>
    </div>
    <div id="list">
        <div id="optimistic" class="text-primary" data-show="$fetching" data-text="$name"></div>
    </div>
  </div>
</div>
