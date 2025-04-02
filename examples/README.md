# Suasor Examples

This directory contains example applications showcasing how to use various Suasor client capabilities.

## Available Examples

### Claude AI Client

`claude_client_example.go` - A simple example showing how to use the Claude AI client for text generation.

To run:
```bash
CLAUDE_API_KEY=your-api-key make claude-example
```

### Movie Recommendations

`movie_recommendations.go` - A comprehensive example combining:

1. Emby media client to fetch user's watch history and favorites
2. Claude AI client to generate personalized movie recommendations
3. Checking if recommended movies already exist in the user's library

This example demonstrates how to use multiple clients together to create a more powerful application.

To run:
```bash
make movie-recommendations
```

This example will use the credentials from your `.env` file. Ensure it contains:
- `EMBY_TEST_URL`
- `EMBY_TEST_API_KEY`
- `EMBY_TEST_USER`
- `CLAUDE_API_KEY`

#### How it works

1. The application loads environment variables from the `.env` file
2. It connects to your Emby server and retrieves recently watched and favorite movies
3. It sends this information to Claude AI to generate personalized recommendations
4. It checks if the recommended movies already exist in your library
5. It saves the recommendations as a JSON file

#### Sample output

```json
{
  "recommendations": [
    {
      "title": "The Shawshank Redemption",
      "year": 1994,
      "reasons": [
        "Strong character-driven narrative similar to your favorites",
        "Themes of hope and resilience"
      ],
      "genreMatch": ["Drama"],
      "directorRef": "Directed by Frank Darabont who has a similar style to directors you enjoy"
    }
  ],
  "basedOn": [
    "The Godfather",
    "Pulp Fiction"
  ],
  "generatedAt": "2025-04-01T12:30:45Z"
}
```

## Creating New Examples

To create a new example:

1. Create a new Go file in the `examples` directory
2. Add appropriate imports for the clients you want to showcase
3. Add a new target in the Makefile to run your example
4. Update this README to document your example