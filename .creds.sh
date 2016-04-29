export VCAP_SERVICES='{
    "alchemy_api": [
       {
          "name": "alchemyAPI",
          "label": "alchemy_api",
          "plan": "ecosystem",
          "credentials": {
             "url": "https://gateway-a.watsonplatform.net/calls",
             "apikey": "",
          }
       }
    ],
    "concept_insights": [
	{
	    "name": "concept-insights-service",
	    "label": "concept_insights",
	    "plan": "standard",
	    "credentials": {
		"url": "https://gateway.watsonplatform.net/concept-insights/api",
		"username": "",
		"password": ""
	    }
	}
    ],
    "conversations": [
        {
	    "name": "conversations-service",
	    "label": "conversations",
	    "plan": "beta",
	    "credentials": {
		"url": "https://gateway-s.watsonplatform.net/dialog-beta/api",
		"password": "",
		"username": ""
	    }
	}
    ],
    "dialog": [
        {
	    "name": "dialog-service",
	    "label": "dialog",
	    "plan": "standard",
	    "credentials": {
		"url": "https://gateway.watsonplatform.net/dialog/api",
		"username": "",
		"password": ""
	    }
	}
    ],
    "document_conversion": [
       {
          "name": "document-conversion",
          "label": "document_conversion",
          "plan": "enterprise",
          "credentials": {
             "url": "https://gateway.watsonplatform.net/document-conversion/api",
             "password": "",
             "username": ""
          }
       }
    ],
    "language_translation": [
       {
          "name": "languague-translation",
          "label": "language_translation",
          "plan": "standard_customizable",
          "credentials": {
             "url": "https://gateway.watsonplatform.net/language-translation/api",
             "password": "",
             "username": ""
          }
       }
    ],
    "natural_language_classifier": [
	{
	    "name": "classifier-service",
	    "label": "natural_language_classifier",
	    "plan": "standard",
	    "credentials": {
		"url": "https://gateway.watsonplatform.net/natural-language-classifier/api",
		"username": "",
		"password": ""
	    }
	}
    ],
    "personality_insights" : [
	{
	    "name": "personality-insights",
	    "label": "personality_insights",
	    "plan": "ecosystem",
	    "credentials": {
	 	"url": "https://gateway.watsonplatform.net/personality-insights/api",
		"password": "",
		"username": ""
	    }
	}
    ],
    "retrieve_and_rank": [
       {
          "name": "retrieve-and-rank",
          "label": "retrieve_and_rank",
          "plan": "enterprise",
          "credentials": {
             "url": "https://gateway.watsonplatform.net/retrieve-and-rank/api",
             "password": "",
             "username": ""
          }
       }
    ],
    "speech_to_text": [
	{
	    "name": "speech-to-text-service",
	    "label": "speech_to_text",
	    "plan": "standard",
	    "credentials": {
		"url": "https://stream.watsonplatform.net/speech-to-text/api",
		"username": "",
		"password": ""
	    }
	}
    ],
    "text_to_speech" : [
	{
	    "name": "text-to-speech",
	    "label": "text_to_speech",
	    "credentials": {
		"url": "https://stream.watsonplatform.net/text-to-speech/api",
		"password": "",
		"username": ""
	    }
	}
    ],
    "tone_analyzer": [
       {
          "name": "tone-analyzer",
          "label": "tone_analyzer",
          "plan": "beta",
          "credentials": {
             "url": "https://gateway.watsonplatform.net/tone-analyzer-beta/api",
             "password": "",
             "username": ""
          }
       }
    ],
    "tradeoff_analytics": [
       {
          "name": "tradeoff-analytics",
          "label": "tradeoff_analytics",
          "plan": "enterprise",
          "credentials": {
             "url": "https://gateway.watsonplatform.net/tradeoff-analytics/api",
             "password": "",
             "username": ""
          }
       }
    ],
    "visual_insights": [
       {
          "name": "visual-insights",
          "label": "visual_insights",
          "plan": "experimental",
          "credentials": {
             "url": "https://gateway.watsonplatform.net/visual-insights-experimental/api",
             "password": "",
             "username": ""
         }
       }
    ],
    "visual_recognition": [
       {
          "name": "visual-recognition",
          "label": "visual_recognition",
          "plan": "free",
          "credentials": {
             "url": "https://gateway.watsonplatform.net/visual-recognition-beta/api",
             "password": "",
             "username": ""
          }
       }
    ]
}'
