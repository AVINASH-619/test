const express = require('express');
const multer = require('multer');
const fetch = require('node-fetch');

const app = express();
const port = 3000;

app.use(express.static('public'));

const storage = multer.memoryStorage();
const upload = multer({ storage: storage });

// Replace 'YOUR_GITHUB_TOKEN', 'USERNAME', 'REPOSITORY', and 'BRANCH' with your GitHub token, username, repository name, and branch name.
const githubToken = 'YOUR_GITHUB_TOKEN';
const githubUsername = 'USERNAME';
const githubRepo = 'REPOSITORY';
const githubBranch = 'BRANCH';

app.post('/upload', upload.single('file'), async (req, res) => {
    const fileBuffer = req.file.buffer;

    const uploadUrl = `https://api.github.com/repos/${githubUsername}/${githubRepo}/contents/temp_file.txt`;
    
    try {
        const response = await fetch(uploadUrl, {
            method: 'PUT',
            headers: {
                'Authorization': `Bearer ${githubToken}`,
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                message: 'Add file',
                content: fileBuffer.toString('base64'),
                branch: githubBranch,
            }),
        });

        if (response.ok) {
            console.log('File uploaded successfully!');
            res.sendStatus(200);
        } else {
            const errorMessage = await response.text();
            console.error(`Failed to upload file. ${errorMessage}`);
            res.status(500).send('Internal Server Error');
        }
    } catch (error) {
        console.error('Error:', error);
        res.status(500).send('Internal Server Error');
    }
});

app.listen(port, () => {
    console.log(`Server is running at http://localhost:${port}`);
});
