import { useState } from 'react'
import './App.css'

function App() {
  const [selectedFile, setSelectedFile] = useState(null);
  const [previewUrl, setPreviewUrl] = useState('');
  const [isProcessing, setIsProcessing] = useState(false);
  const [downloadUrl, setDownloadUrl] = useState('');
  const [statusMessage, setStatusMessage] = useState('');
  const [uploadedImagePath, setUploadedImagePath] = useState('');

  const handleFileChange = (event) => {
    const file = event.target.files[0];
    if (file && file.type === 'image/png') {
      setSelectedFile(file);
      setPreviewUrl(URL.createObjectURL(file));
      setStatusMessage(`已选择文件: ${file.name}`);
      
      // 上传文件到后端
      uploadFile(file);
    } else {
      setStatusMessage('请选择PNG格式的图片文件');
    }
  };

  const handleDragOver = (event) => {
    event.preventDefault();
  };

  const handleDrop = (event) => {
    event.preventDefault();
    const file = event.dataTransfer.files[0];
    if (file && file.type === 'image/png') {
      setSelectedFile(file);
      setPreviewUrl(URL.createObjectURL(file));
      setStatusMessage(`已选择文件: ${file.name}`);
      
      // 上传文件到后端
      uploadFile(file);
    } else {
      setStatusMessage('请选择PNG格式的图片文件');
    }
  };

  const uploadFile = async (file) => {
    const formData = new FormData();
    formData.append('image', file);
    
    try {
      setStatusMessage('正在上传图片...');
      const response = await fetch('http://localhost:8080/api/v1/upload', {
        method: 'POST',
        body: formData,
      });
      
      if (response.ok) {
        const data = await response.json();
        setUploadedImagePath(data.filename); // 后端返回的是filename字段
        setStatusMessage('图片上传成功');
      } else {
        setStatusMessage('图片上传失败');
      }
    } catch (error) {
      console.error('上传错误:', error);
      setStatusMessage('图片上传过程中发生错误');
    }
  };

  const handleProcess = async () => {
    if (!uploadedImagePath) {
      setStatusMessage('请先选择并上传图片文件');
      return;
    }

    setIsProcessing(true);
    setStatusMessage('正在处理图片...');
    
    try {
      const response = await fetch('http://localhost:8080/api/v1/process', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ filename: uploadedImagePath }),
      });
      
      if (response.ok) {
        const data = await response.json();
        setDownloadUrl(`http://localhost:8080${data.download_url}`); // 后端返回的是相对路径，需要拼接完整URL
        setStatusMessage('图片处理完成，可以下载结果');
      } else {
        setStatusMessage('图片处理失败');
      }
    } catch (error) {
      console.error('处理错误:', error);
      setStatusMessage('图片处理过程中发生错误');
    } finally {
      setIsProcessing(false);
    }
  };

  const handleDownload = () => {
    if (downloadUrl) {
      setStatusMessage('开始下载...');
      // 使用window.location.href触发下载
      window.location.href = downloadUrl;
      setStatusMessage('下载完成');
    }
  };

  return (
    <div className="app">
      <header className="app-header">
        <h1>图集切割工具</h1>
      </header>
      
      <main className="app-main">
        <section className="upload-section">
          <h2>上传图片</h2>
          <div 
            className="drop-area"
            onDragOver={handleDragOver}
            onDrop={handleDrop}
          >
            <p>拖拽PNG图片到此处或</p>
            <input 
              type="file" 
              accept="image/png" 
              onChange={handleFileChange} 
              id="file-input"
              style={{display: 'none'}}
            />
            <label htmlFor="file-input" className="file-input-label">
              选择文件
            </label>
          </div>
          {selectedFile && (
            <p className="file-info">已选择: {selectedFile.name}</p>
          )}
        </section>
        
        <section className="preview-section">
          <h2>预览</h2>
          {previewUrl && (
            <div className="image-preview">
              <img src={previewUrl} alt="预览" />
            </div>
          )}
        </section>
        
        <section className="actions-section">
          <h2>操作</h2>
          <button 
            onClick={handleProcess} 
            disabled={isProcessing || !uploadedImagePath}
            className="process-button"
          >
            {isProcessing ? '处理中...' : '开始切割'}
          </button>
          {downloadUrl && (
            <button 
              onClick={handleDownload} 
              className="download-button"
            >
              下载结果
            </button>
          )}
        </section>
      </main>
      
      <footer className="app-footer">
        <div className="status-message">
          {statusMessage}
        </div>
      </footer>
    </div>
  )
}

export default App
