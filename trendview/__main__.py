

from ollama import chat
from ollama import ChatResponse


def main():
    response: ChatResponse = chat(model='mistral', messages=[
    {
        'role': 'user',
        'content': 'Why is the sky blue?',
    },
    ])
    print(response['message']['content'])
    # or access fields directly from the response object
    print(response.message.content)

if (__name__ == "__main__"):
    main()
